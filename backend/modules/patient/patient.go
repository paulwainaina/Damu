package patient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Patient struct {
	ID          uint64 `bson:"ID"`
	Name        string `bson:"Name"`
	DateofBirth string `bson:"DateofBirth"`
	Gender      string `bson:"Gender"`
	BloodGroup  string `bson:"BloodGroup"`
}

func (patient *Patient) UnmarshalJSON(data []byte) error {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}
	for k, v := range jsonData {
		switch strings.ToLower(k) {
		case "id":{
			patient.ID=uint64(v.(float64))
		}
		case "name":{
			patient.Name=v.(string)
		}
		case "gender":{
			patient.Gender=v.(string)
		}
		case "bloodgroup":{
			patient.BloodGroup=v.(string)
		}
		case "dateofbirth":{			
			_,err:=time.Parse(time.DateOnly,v.(string))
			if err!=nil{
				return err
			}
			patient.DateofBirth=v.(string)
		}
		}
	}
	return nil
}

type Patients struct{
	Patients []*Patient	
}

func NewPatients()*Patients{
	return &Patients{Patients: make([]*Patient, 0)}
}
func (patients *Patients) GenerateNewID() uint64 {
	var x uint64 = 0
	if len(patients.Patients) == 0 {
		return 1
	}
	var ids []uint64 = make([]uint64, 0)
	for _, i := range patients.Patients {
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

func (patients *Patients)GetPatient(id uint64)(*Patient,error){
	for _,p:=range patients.Patients{
		if p.ID==id{
			return p,nil
		}
	}
	return &Patient{},fmt.Errorf("patient %v not found",id)
}
func (patients *Patients)AddPatient(p Patient)(*Patient,error){
	if p.ID!=0{
		return &Patient{},fmt.Errorf("new patient cannot have id %v",p.ID)
	}
	p.ID=patients.GenerateNewID()
	patients.Patients=append(patients.Patients, &p)
	return &p,nil
}
func (patients *Patients)DeletePatient(id uint64)(*Patient,error){
	for i,p:=range patients.Patients{
		if p.ID==id{
			patients.Patients=append(patients.Patients[:i],patients.Patients[i+1:]... )
			return p,nil
		}
	}
	return &Patient{},fmt.Errorf("patient %v not found",id)
}

func (patients *Patients)UpdatePatient(patient Patient)(*Patient,error){
	for _,p:=range patients.Patients{
		if p.ID==patient.ID{
			*p=patient
			return p,nil
		}
	}
	return &Patient{},fmt.Errorf("patient %v not found",patient.ID)
}


func (patients *Patients)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type Error struct{Error string}
	if r.URL.Path=="/patient"{
		switch(r.Method){
			case http.MethodGet:{
				patient:=Patient{}
				err := json.NewDecoder(r.Body).Decode(&patient)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				p,err:=patients.GetPatient(patient.ID)
				if err!=nil{
					json.NewEncoder(w).Encode(Error{Error: err.Error()})
					return
				}
				json.NewEncoder(w).Encode(p)
				return
			}
			case http.MethodPut:{
				patient:=Patient{}
				err := json.NewDecoder(r.Body).Decode(&patient)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				p,err:=patients.UpdatePatient(patient)
				if err!=nil{
					json.NewEncoder(w).Encode(Error{Error: err.Error()})
					return
				}
				json.NewEncoder(w).Encode(p)
				return
			}
			case http.MethodDelete:{
				patient:=Patient{}
				err := json.NewDecoder(r.Body).Decode(&patient)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				p,err:=patients.DeletePatient(patient.ID)
				if err!=nil{
					json.NewEncoder(w).Encode(Error{Error: err.Error()})
					return
				}
				json.NewEncoder(w).Encode(p)
				return
			}
		}

	}else if (r.URL.Path=="/patients"){
		if r.Method==http.MethodGet{
			json.NewEncoder(w).Encode(patients.Patients)
			return
		}
	}
	w.WriteHeader(http.StatusNotImplemented)
}