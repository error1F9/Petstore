package controller

import (
	"Petstore/models/pet/entity"
	"Petstore/models/pet/service"
	"Petstore/responder"
	"github.com/go-chi/chi/v5"
	"github.com/ptflp/godecoder"
	"net/http"
	"strconv"
)

type PetController interface {
	Add(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	FindByStatus(http.ResponseWriter, *http.Request)
	FindById(http.ResponseWriter, *http.Request)
	UpdateById(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type PetControl struct {
	service service.PetServicer
	responder.Responder
	godecoder.Decoder
}

func NewPetController(service service.PetServicer, responder responder.Responder, decoder godecoder.Decoder) *PetControl {
	return &PetControl{service: service, Responder: responder, Decoder: decoder}
}

// Add Adding a pet to the store
// @Summary Add a pet
// @Description Add a new pet to the store
// @Tags pet
// @Accept  json
// @Produce  json
// @Param   pet body entity.Pet true "Pet object that needs to be added to the store"
// @Success 200 {object} PetResponse
// @Router /pet [post]
func (p *PetControl) Add(w http.ResponseWriter, r *http.Request) {
	var pet entity.Pet
	if err := p.Decode(r.Body, &pet); err != nil {
		p.ErrorBadRequest(w, err)
		return
	}
	out := p.service.Add(r.Context(), service.PetAddIn{
		Name:     pet.Name,
		Category: pet.Category,
		Status:   pet.Status,
	})

	if out.Err != nil {
		p.ErrorInternal(w, out.Err)
		return
	}

	p.OutputJSON(w, PetResponse{
		Code:    http.StatusCreated,
		Success: true,
	})
}

// Update Updating an existing pet
// @Summary Update a pet
// @Description Update an existing pet
// @Tags pet
// @Accept  json
// @Produce  json
// @Param   pet body entity.Pet true "Pet object that needs to be updated"
// @Success 200 {object} PetResponse
// @Router /pet [put]
func (p *PetControl) Update(w http.ResponseWriter, r *http.Request) {
	var pet entity.Pet
	if err := p.Decode(r.Body, &pet); err != nil {
		p.ErrorBadRequest(w, err)
		return
	}
	out := p.service.Update(r.Context(), service.PetUpdateIn{
		Pet: pet,
	})
	if out.Err != nil {
		p.ErrorInternal(w, out.Err)
		return
	}

	p.OutputJSON(w, PetResponse{
		Code:    http.StatusOK,
		Success: true,
	})
}

// FindByStatus Finding pets by their status
// @Summary Find pets by status
// @Description Multiple status values can be provided with comma separated strings
// @Tags pet
// @Accept  json
// @Produce  json
// @Param   status query entity.PetStatus false "Status values that need to be considered for filter"
// @Success 200 {object} PetFindByStatusResponse
// @Router /pet/findByStatus [get]
func (p *PetControl) FindByStatus(w http.ResponseWriter, r *http.Request) {

	str := r.URL.Query().Get("status")
	out := p.service.FindByStatus(r.Context(), service.PetFindByStatusIn{
		Status: entity.PetStatus(str),
	})
	if out.Err != nil {
		p.ErrorInternal(w, out.Err)
		return
	}
	p.OutputJSON(w, PetFindByStatusResponse{
		Code:    http.StatusOK,
		Success: true,
		Data:    DataPets{Pets: out.Pets},
	})
}

// FindById Returning a single pet
// @Summary Find pet by ID
// @Description Returns a single pet
// @Tags pet
// @Accept  json
// @Produce  json
// @Param   id path int true "ID of pet to return"
// @Success 200 {object} PetResponseData
// @Router /pet/{id} [get]
func (p *PetControl) FindById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, _ := strconv.Atoi(id)

	out := p.service.FindById(r.Context(), service.PetFindByIdIn{
		PetID: uint64(intID),
	})
	if out.Err != nil {
		p.ErrorInternal(w, out.Err)
		return
	}

	p.OutputJSON(w, PetResponseData{
		Code:    http.StatusOK,
		Success: true,
		Data:    Data{Pet: out.Pet},
	})

}

// UpdateById Updating an existing pet by ID
// @Summary Update a pet by ID
// @Description Updates pet in the store with form data
// @Tags pet
// @Accept  json
// @Produce  json
// @Param   id path int true "ID of pet that needs to be updated"
// @Param   pet body entity.Pet true "Pet object that needs to be updated"
// @Success 200 {object} PetResponse
// @Router /pet/{id} [post]
func (p *PetControl) UpdateById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, _ := strconv.Atoi(id)

	var pet entity.Pet
	if err := p.Decode(r.Body, &pet); err != nil {
		p.ErrorBadRequest(w, err)
		return
	}

	pet.ID = uint64(intID)

	out := p.service.UpdateById(r.Context(), service.PetUpdateByIdIn{
		Pet: pet,
	})
	if out.Err != nil {
		p.ErrorInternal(w, out.Err)
		return
	}

	p.OutputJSON(w, PetResponse{
		Code:    http.StatusOK,
		Success: true,
	})
}

// Delete Deleting a pet by ID
// @Summary Delete a pet by ID
// @Description Deletes a pet by ID
// @Tags pet
// @Accept  json
// @Produce  json
// @Param   id path int true "ID of pet to delete"
// @Success 200 {object} PetResponse
// @Router /pet/{id} [delete]
func (p *PetControl) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	intID, _ := strconv.Atoi(id)

	out := p.service.Delete(r.Context(), service.PetDeleteIn{
		PetID: uint64(intID),
	})

	if out.Err != nil {
		p.ErrorInternal(w, out.Err)
		return
	}

	p.OutputJSON(w, PetResponse{
		Code:    http.StatusOK,
		Success: true,
	})

}
