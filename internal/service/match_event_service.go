package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/xyz-corp/xyz-football-api/internal/dto"
	"github.com/xyz-corp/xyz-football-api/internal/entity"
	"github.com/xyz-corp/xyz-football-api/internal/repository"
	"github.com/xyz-corp/xyz-football-api/pkg/apperror"
	"gorm.io/gorm"
)

type MatchEventService interface {
	SubmitResult(matchUUID uuid.UUID, req dto.MatchResultRequest, isCorrection bool) (*dto.MatchResponse, error)
	AddEvent(matchUUID uuid.UUID, req dto.MatchEventRequest) (*dto.MatchEventResponse, error)
	ListEvents(matchUUID uuid.UUID) ([]dto.MatchEventResponse, error)
	DeleteEvent(matchUUID uuid.UUID, eventUUID uuid.UUID) error
}

type matchEventService struct {
	matchRepo      repository.MatchRepository
	playerRepo     repository.PlayerRepository
	matchEventRepo repository.MatchEventRepository
}

func NewMatchEventService(matchRepo repository.MatchRepository, playerRepo repository.PlayerRepository, matchEventRepo repository.MatchEventRepository) MatchEventService {
	return &matchEventService{matchRepo: matchRepo, playerRepo: playerRepo, matchEventRepo: matchEventRepo}
}

// validatedEvent holds a request event resolved against its DB player record.
type validatedEvent struct {
	req    dto.MatchEventRequest
	player *entity.Player
	assist *entity.Player
}

// resolveAndValidateEvents applies rules FR-4.3 and TRD 8.7 items 3-6 against a batch of
// requested events for a given match. It does not touch the database beyond player lookups.
func (s *matchEventService) resolveAndValidateEvents(match *entity.Match, reqs []dto.MatchEventRequest) ([]validatedEvent, error) {
	resolved := make([]validatedEvent, 0, len(reqs))

	seen := map[string]bool{}
	redCardCount := map[uint64]int{}
	yellowCardMinutes := map[uint64][]int16{}

	for _, r := range reqs {
		player, err := s.playerRepo.FindByUUID(r.PlayerUUID)
		if err != nil {
			return nil, wrapNotFound(err, "Pemain "+r.PlayerUUID.String()+" tidak ditemukan")
		}
		if player.TeamID != match.HomeTeamID && player.TeamID != match.AwayTeamID {
			return nil, apperror.BusinessRule("Pemain " + player.Name + " bukan peserta pertandingan ini")
		}

		var assist *entity.Player
		if r.AssistPlayerUUID != nil {
			if r.EventType != entity.EventTypeGoal {
				return nil, apperror.BusinessRule("assist_player_uuid hanya berlaku untuk event_type GOAL")
			}
			assist, err = s.playerRepo.FindByUUID(*r.AssistPlayerUUID)
			if err != nil {
				return nil, wrapNotFound(err, "Pemain assist tidak ditemukan")
			}
			if assist.TeamID != player.TeamID {
				return nil, apperror.BusinessRule("Pemain assist harus satu tim dengan pencetak gol")
			}
			if assist.ID == player.ID {
				return nil, apperror.BusinessRule("Pemain assist tidak boleh sama dengan pencetak gol")
			}
		}

		key := r.PlayerUUID.String() + "|" + r.EventType + "|" + itoa(r.Minute)
		if seen[key] {
			return nil, apperror.BusinessRule("Event duplikat: pemain, tipe event, dan menit yang sama sudah tercatat")
		}
		seen[key] = true

		if r.EventType == entity.EventTypeRedCard {
			redCardCount[player.ID]++
			if redCardCount[player.ID] > 1 {
				return nil, apperror.BusinessRule("Pemain " + player.Name + " tidak boleh menerima lebih dari satu kartu merah")
			}
		}
		if r.EventType == entity.EventTypeYellowCard {
			yellowCardMinutes[player.ID] = append(yellowCardMinutes[player.ID], r.Minute)
		}

		resolved = append(resolved, validatedEvent{req: r, player: player, assist: assist})
	}

	// Rule 6: seorang pemain hanya boleh dapat kartu kuning kedua bila disertai kartu merah di menit yang sama.
	for playerID, minutes := range yellowCardMinutes {
		if len(minutes) <= 1 {
			continue
		}
		if len(minutes) > 2 {
			return nil, apperror.BusinessRule("Pemain tidak dapat menerima lebih dari dua kartu kuning dalam satu pertandingan")
		}
		if redCardCount[playerID] < 1 {
			return nil, apperror.BusinessRule("Kartu kuning kedua untuk seorang pemain wajib disertai kartu merah pada menit yang sama")
		}
	}

	return resolved, nil
}

func itoa(n int16) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [6]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func (s *matchEventService) SubmitResult(matchUUID uuid.UUID, req dto.MatchResultRequest, isCorrection bool) (*dto.MatchResponse, error) {
	match, err := s.matchRepo.FindByUUID(matchUUID)
	if err != nil {
		return nil, wrapNotFound(err, "Jadwal pertandingan tidak ditemukan")
	}

	if !isCorrection {
		if match.Status == entity.MatchStatusFinished {
			return nil, apperror.BusinessRule("Pertandingan sudah memiliki hasil akhir (FR-4.2)")
		}
		if match.Status == entity.MatchStatusPostpone {
			return nil, apperror.BusinessRule("Pertandingan berstatus postpone, tidak dapat dilaporkan hasilnya")
		}
	} else if match.Status != entity.MatchStatusFinished {
		return nil, apperror.BusinessRule("Koreksi hasil hanya berlaku untuk pertandingan yang sudah finished")
	}

	if match.MatchDate.After(time.Now()) {
		return nil, apperror.BusinessRule("Tanggal pertandingan tidak boleh di masa depan")
	}

	resolved, err := s.resolveAndValidateEvents(match, req.Events)
	if err != nil {
		return nil, err
	}

	goalCount, homeGoals, awayGoals := 0, 0, 0
	for _, ve := range resolved {
		if ve.req.EventType != entity.EventTypeGoal {
			continue
		}
		goalCount++
		if ve.player.TeamID == match.HomeTeamID {
			homeGoals++
		} else {
			awayGoals++
		}
	}
	if int16(goalCount) > req.HomeScore+req.AwayScore {
		return nil, apperror.BusinessRule("Jumlah event GOAL tidak boleh melebihi total skor")
	}
	if int16(homeGoals) > req.HomeScore {
		return nil, apperror.BusinessRule("Jumlah gol tim home melebihi skor home yang dilaporkan")
	}
	if int16(awayGoals) > req.AwayScore {
		return nil, apperror.BusinessRule("Jumlah gol tim away melebihi skor away yang dilaporkan")
	}

	txErr := s.matchRepo.DB().Transaction(func(tx *gorm.DB) error {
		if isCorrection {
			if err := s.matchEventRepo.SoftDeleteByMatchID(tx, match.ID); err != nil {
				return err
			}
		}

		events := make([]entity.MatchEvent, 0, len(resolved))
		for _, ve := range resolved {
			var assistID *uint64
			if ve.assist != nil {
				id := ve.assist.ID
				assistID = &id
			}
			events = append(events, entity.MatchEvent{
				MatchID:        match.ID,
				PlayerID:       ve.player.ID,
				AssistPlayerID: assistID,
				EventType:      ve.req.EventType,
				Minute:         ve.req.Minute,
			})
		}
		if err := s.matchEventRepo.CreateMany(tx, events); err != nil {
			return err
		}

		match.HomeScore = req.HomeScore
		match.AwayScore = req.AwayScore
		match.Status = entity.MatchStatusFinished
		return tx.Save(match).Error
	})
	if txErr != nil {
		return nil, apperror.Internal(txErr.Error())
	}

	updated, err := s.matchRepo.FindByUUIDWithEvents(matchUUID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	return toMatchResponseWithEvents(updated), nil
}

func (s *matchEventService) AddEvent(matchUUID uuid.UUID, req dto.MatchEventRequest) (*dto.MatchEventResponse, error) {
	match, err := s.matchRepo.FindByUUID(matchUUID)
	if err != nil {
		return nil, wrapNotFound(err, "Jadwal pertandingan tidak ditemukan")
	}
	if match.Status != entity.MatchStatusInProgress {
		return nil, apperror.BusinessRule("Event hanya dapat ditambahkan saat pertandingan berstatus in_progress")
	}

	resolved, err := s.resolveAndValidateEvents(match, []dto.MatchEventRequest{req})
	if err != nil {
		return nil, err
	}
	ve := resolved[0]

	dup, err := s.matchEventRepo.ExistsDuplicate(match.ID, ve.player.ID, req.EventType, req.Minute)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}
	if dup {
		return nil, apperror.BusinessRule("Event duplikat: pemain, tipe event, dan menit yang sama sudah tercatat")
	}

	var assistID *uint64
	if ve.assist != nil {
		id := ve.assist.ID
		assistID = &id
	}

	event := entity.MatchEvent{
		MatchID:        match.ID,
		PlayerID:       ve.player.ID,
		AssistPlayerID: assistID,
		EventType:      req.EventType,
		Minute:         req.Minute,
	}
	if err := s.matchEventRepo.CreateMany(s.matchRepo.DB(), []entity.MatchEvent{event}); err != nil {
		return nil, apperror.Internal(err.Error())
	}

	resp := toMatchEventResponse(&event, match)
	return &resp, nil
}

func (s *matchEventService) ListEvents(matchUUID uuid.UUID) ([]dto.MatchEventResponse, error) {
	match, err := s.matchRepo.FindByUUID(matchUUID)
	if err != nil {
		return nil, wrapNotFound(err, "Jadwal pertandingan tidak ditemukan")
	}

	events, err := s.matchEventRepo.ListByMatch(match.ID)
	if err != nil {
		return nil, apperror.Internal(err.Error())
	}

	result := make([]dto.MatchEventResponse, 0, len(events))
	for i := range events {
		result = append(result, toMatchEventResponse(&events[i], match))
	}
	return result, nil
}

func (s *matchEventService) DeleteEvent(matchUUID uuid.UUID, eventUUID uuid.UUID) error {
	match, err := s.matchRepo.FindByUUID(matchUUID)
	if err != nil {
		return wrapNotFound(err, "Jadwal pertandingan tidak ditemukan")
	}

	event, err := s.matchEventRepo.FindByUUID(match.ID, eventUUID)
	if err != nil {
		return wrapNotFound(err, "Event tidak ditemukan")
	}

	return s.matchEventRepo.Delete(event)
}
