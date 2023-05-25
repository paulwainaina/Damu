package triage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Triage struct {
	ID          uint64 `bson:"ID"`
	PatientID   uint64 `bson:"PatientID"`
	DateofVisit string `bson:"DateofVisit"`
	Details     string `bson:"Details"`
	DoctorName  string `bson:"DoctorName"`
}

func (triage *Triage) UnmarshalJSON(data []byte) error {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}
	for k, v := range jsonData {
		switch strings.ToLower(k) {
		case "id":
			{
				if v.(string)==""{ break}
				x,err:=strconv.ParseUint(v.(string),10,64)
				if err != nil {
					return err
				}
				triage.ID = x
			}
		case "doctorname":
			{
				triage.DoctorName = v.(string)
			}
		case "details":
			{
				triage.Details = v.(string)
			}
		case "patientid":
			{
				triage.PatientID = uint64(v.(float64))
			}
		case "dateofvisit":
			{
				_, err := time.Parse(time.DateOnly, v.(string))
				if err != nil {
					return err
				}
				triage.DateofVisit = v.(string)
			}
		}
	}
	return nil
}

type Triages struct {
	Triages []*Triage
	pattern *regexp.Regexp
}

func NewTriages() *Triages {
	return &Triages{Triages: make([]*Triage, 0), pattern: regexp.MustCompile(`^/triage/(\d+)/?`)}
}
func (triages *Triages) GenerateNewID() uint64 {
	var x uint64 = 0
	if len(triages.Triages) == 0 {
		return 1
	}
	var ids []uint64 = make([]uint64, 0)
	for _, i := range triages.Triages {
		ids = append(ids, i.ID)
	}
	sort.Slice(ids, func(a, b int) bool { return ids[a] < ids[b] })
	for i := 1; i < len(ids); i++ {
		if ids[i-1]+1 != ids[i] {
			x = ids[i-1] + 1
			break
		}
	}
	if x == 0 {
		x = ids[len(ids)-1] + 1
	}
	return x
}

func (triages *Triages) GetTriage(id uint64) (*Triage, error) {
	for _, p := range triages.Triages {
		if p.ID == id {
			return p, nil
		}
	}
	return &Triage{}, fmt.Errorf("triage %v not found", id)
}
func (triages *Triages) GetTriagePatient(id uint64) (*Triage, error) {
	for _, p := range triages.Triages {
		if p.PatientID == id {
			return p, nil
		}
	}
	return &Triage{}, fmt.Errorf("triage %v not found", id)
}
func (triages *Triages) AddTriage(p Triage) (*Triage, error) {
	if p.ID != 0 {
		return &Triage{}, fmt.Errorf("new triage cannot have id %v", p.ID)
	}
	p.ID = triages.GenerateNewID()
	triages.Triages = append(triages.Triages, &p)
	return &p, nil
}
func (triages *Triages) DeleteTriage(id uint64) (*Triage, error) {
	for i, p := range triages.Triages {
		if p.ID == id {
			triages.Triages = append(triages.Triages[:i], triages.Triages[i+1:]...)
			return p, nil
		}
	}
	return &Triage{}, fmt.Errorf("triage %v not found", id)
}

func (triages *Triages) UpdateTriage(triage Triage) (*Triage, error) {
	for _, p := range triages.Triages {
		if p.ID == triage.ID {
			*p = triage
			return p, nil
		}
	}
	return &Triage{}, fmt.Errorf("triage %v not found", triage.ID)
}

func (triages *Triages) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	type Error struct{ Error string }
	if r.URL.Path == "/triage" {
		switch r.Method {
		
		case http.MethodPut:
			{
				triage := Triage{}
				err := json.NewDecoder(r.Body).Decode(&triage)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				p, err := triages.UpdateTriage(triage)
				if err != nil {
					json.NewEncoder(w).Encode(Error{Error: err.Error()})
					return
				}
				json.NewEncoder(w).Encode(p)
				return
			}
		case http.MethodPost:
			{
				triage := Triage{}
				err := json.NewDecoder(r.Body).Decode(&triage)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				p, err := triages.AddTriage(triage)
				if err != nil {
					json.NewEncoder(w).Encode(Error{Error: err.Error()})
					return
				}
				json.NewEncoder(w).Encode(p)
				return
			}
		case http.MethodDelete:
			{
				triage := Triage{}
				err := json.NewDecoder(r.Body).Decode(&triage)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				p, err := triages.DeleteTriage(triage.ID)
				if err != nil {
					json.NewEncoder(w).Encode(Error{Error: err.Error()})
					return
				}
				json.NewEncoder(w).Encode(p)
				return
			}
		}

	} else if r.URL.Path == "/triages" {
		if r.Method == http.MethodGet {
			json.NewEncoder(w).Encode(triages.Triages)
			return
		}
	}else{
		matches :=triages.pattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		switch(r.Method){
			case http.MethodGet:
			{
				
				p, err := triages.GetTriage(uint64(id))
				if err != nil {
					json.NewEncoder(w).Encode(Error{Error: err.Error()})
					return
				}
				json.NewEncoder(w).Encode(p)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotImplemented)
}
