package services

import (
	"fmt"
	"food/src/api/models/ingredient"
	"food/src/api/models/receipt"
	"food/src/api/models/tools"
	"github.com/jinzhu/gorm"
	"mime/multipart"
	"strings"
)

func GetReceiptService(db *gorm.DB) *Receipt {
	return &Receipt{
		receiptRepo: receipt.GetReceiptRepository(db),
		ingredientRepo: ingredient.GetMediaRepository(db),
		mediaSvc: GetMediaService(db),
	}
}

type Receipt struct {
	receiptRepo *receipt.ReceiptRepository
	ingredientRepo *ingredient.IngredientRepository
	mediaSvc *Media
}

func (s *Receipt) GetAllReceiptsByCategory(category string) (receipts []receipt.Receipt, err error) {
	receipts, err = s.receiptRepo.GetAllByCategory(category)
	return
}

func (s *Receipt) GetAllReceiptIngredientsById(id uint) (ingredients []receipt.ReceiptIngredient, err error) {
	_, err = s.receiptRepo.GetById(id)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}
	if err != nil {
		return
	}
	ingredients, err = s.receiptRepo.GetIngredientsById(id)
	return
}

func (s *Receipt) GetAllReceiptDirectionsById(id uint) (directions []receipt.ReceiptDirection, err error) {
	_, err = s.receiptRepo.GetById(id)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}
	if err != nil {
		return
	}
	directions, err = s.receiptRepo.GetDirectionsById(id)
	return
}

func (s *Receipt) GetAllIngredients() (receipts []ingredient.Ingredient, err error) {
	receipts, err = s.ingredientRepo.GetAll()
	return
}

type CreateIngredientRequest struct {
	// (required)
	Name    string     `json:"name" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
}

func (u *CreateIngredientRequest) TrimSpaces() {
	u.Name = strings.TrimSpace(u.Name)
}

type UpdateIngredientRequest struct {
	CreateIngredientRequest
}

type CreateReceiptRequest struct {
	// (required)
	Name    string     `json:"name" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
	Description    string     `json:"description" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
	Category    string     `json:"category" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
	CookingTime    int     `json:"cooking_time" minimum:"1" binding:"required" validate:"min=1"`
}

func (u *CreateReceiptRequest) TrimSpaces() {
	u.Name = strings.TrimSpace(u.Name)
	u.Description = strings.TrimSpace(u.Description)
	u.Category = strings.TrimSpace(u.Category)
}

type UpdateReceiptRequest struct {
	// (required)
	Name    string     `json:"name" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
	Description    string     `json:"description" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
	Category    string     `json:"category" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
	CookingTime    int     `json:"cooking_time" minimum:"1" binding:"required" validate:"min=1"`
}

func (u *UpdateReceiptRequest) TrimSpaces() {
	u.Name = strings.TrimSpace(u.Name)
	u.Description = strings.TrimSpace(u.Description)
	u.Category = strings.TrimSpace(u.Category)
}

type CreateReceiptIngredientRequest struct {
	// (required)
	Quantity    string     `json:"quantity" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
	IngredientId    uint     `json:"ingredient_id" minimum:"1" binding:"required" validate:"min=1"`
}

func (u *CreateReceiptIngredientRequest) TrimSpaces() {
	u.Quantity = strings.TrimSpace(u.Quantity)
}

type UpdateReceiptIngredientRequest struct {
	// (required)
	Quantity    string     `json:"quantity" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
}

func (u *UpdateReceiptIngredientRequest) TrimSpaces() {
	u.Quantity = strings.TrimSpace(u.Quantity)
}

type CreateReceiptDirectionRequest struct {
	// (required)
	Description    string     `json:"description" minLength:"3" maxLength:"255" binding:"required" validate:"max=255,min=3"`
}

func (u *CreateReceiptDirectionRequest) TrimSpaces() {
	u.Description = strings.TrimSpace(u.Description)
}

type UpdateReceiptDirectionRequest struct {
	CreateReceiptDirectionRequest
}

func (s *Receipt) CreateIngredient(ingredientRequest CreateIngredientRequest) (i ingredient.Ingredient, err error) {
	ingredientRequest.TrimSpaces()
	err = tools.Validator.Struct(ingredientRequest)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}
	i = ingredient.Ingredient{Name:ingredientRequest.Name}
	err = s.ingredientRepo.Create(&i)
	return
}

func (s *Receipt) UpdateIngredient(id uint, ingredientRequest UpdateIngredientRequest) (i ingredient.Ingredient, err error) {
	ingredientRequest.TrimSpaces()
	err = tools.Validator.Struct(ingredientRequest)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}

	oldItem, err := s.ingredientRepo.GetById(id)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}
	if err != nil {
		return
	}
	i = ingredient.Ingredient{Id: oldItem.Id, Name: ingredientRequest.Name, CreatedAt: oldItem.CreatedAt}
	err = s.ingredientRepo.Update(&i)
	return
}

func (s *Receipt) CreateReceipt(userId uint, request CreateReceiptRequest) (i receipt.Receipt, err error) {
	request.TrimSpaces()
	err = tools.Validator.Struct(request)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}
	i = receipt.Receipt{
		Name: request.Name,
		Description: request.Description,
		Category: request.Category,
		CookingTime: request.CookingTime,
		UserId: userId,

	}
	err = s.receiptRepo.Create(&i)
	return
}

func (s *Receipt) CreateReceiptIngredient(receiptId, userId uint, request CreateReceiptIngredientRequest) (i receipt.ReceiptIngredient, err error) {
	request.TrimSpaces()
	err = tools.Validator.Struct(request)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}

	r, err := s.receiptRepo.GetById(receiptId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}

	if err != nil {
		return
	}

	if r.UserId != userId {
		err = tools.NewNotPermittedErr(fmt.Errorf("user id mismached"))
		return
	}

	_, err = s.ingredientRepo.GetById(request.IngredientId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}

	if err != nil {
		return
	}

	i = receipt.ReceiptIngredient{
		Quantity: request.Quantity,
		ReceiptId: receiptId,
		IngredientId: request.IngredientId,
	}
	err = s.receiptRepo.CreateIngredient(&i)
	return
}

func (s *Receipt) UpdateReceiptIngredient(receiptId, rIngredientId, userId uint, request UpdateReceiptIngredientRequest) (i receipt.ReceiptIngredient, err error) {
	request.TrimSpaces()
	err = tools.Validator.Struct(request)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}

	r, err := s.receiptRepo.GetById(receiptId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}

	if err != nil {
		return
	}

	if r.UserId != userId {
		err = tools.NewNotPermittedErr(fmt.Errorf("user id mismached"))
		return
	}

	_, err = s.receiptRepo.GetIngredientById(rIngredientId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}

	if err != nil {
		return
	}

	i = receipt.ReceiptIngredient{
		Id: rIngredientId,
		Quantity: request.Quantity,
	}
	err = s.receiptRepo.UpdateIngredient(&i)
	if err != nil {
		return
	}

	i, err = s.receiptRepo.GetIngredientById(i.Id)
	return
}

func (s *Receipt) DeleteReceiptIngredient(receiptId, rIngredientId, userId uint) (err error) {
	r, err := s.receiptRepo.GetById(receiptId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}

	if err != nil {
		return
	}

	if r.UserId != userId {
		err = tools.NewNotPermittedErr(fmt.Errorf("user id mismached"))
		return
	}

	err = s.receiptRepo.DeleteIngredientById(rIngredientId)
	if gorm.IsRecordNotFoundError(err) {
		err = nil
		return
	}

	if err != nil {
		return
	}
	return
}


// -------
func (s *Receipt) CreateReceiptDirection(receiptId, userId uint, request CreateReceiptDirectionRequest) (i receipt.ReceiptDirection, err error) {
	request.TrimSpaces()
	err = tools.Validator.Struct(request)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}

	r, err := s.receiptRepo.GetById(receiptId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}

	if err != nil {
		return
	}

	if r.UserId != userId {
		err = tools.NewNotPermittedErr(fmt.Errorf("user id mismached"))
		return
	}

	if err != nil {
		return
	}

	i = receipt.ReceiptDirection{
		ReceiptId: receiptId,
		Description: request.Description,
	}
	err = s.receiptRepo.CreateDirection(&i)
	return
}

func (s *Receipt) UpdateReceiptDirection(receiptId, rDirectionId, userId uint, request UpdateReceiptDirectionRequest) (i receipt.ReceiptDirection, err error) {
	request.TrimSpaces()
	err = tools.Validator.Struct(request)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}

	r, err := s.receiptRepo.GetById(receiptId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}

	if err != nil {
		return
	}

	if r.UserId != userId {
		err = tools.NewNotPermittedErr(fmt.Errorf("user id mismached"))
		return
	}

	_, err = s.receiptRepo.GetDirectionById(rDirectionId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}

	if err != nil {
		return
	}

	i = receipt.ReceiptDirection{
		Id:       rDirectionId,
		Description: request.Description,
	}
	err = s.receiptRepo.UpdateDirection(&i)
	if err != nil {
		return
	}

	i, err = s.receiptRepo.GetDirectionById(i.Id)
	return
}

func (s *Receipt) DeleteReceiptDirection(receiptId, rDirectionId, userId uint) (err error) {
	r, err := s.receiptRepo.GetById(receiptId)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}

	if err != nil {
		return
	}

	if r.UserId != userId {
		err = tools.NewNotPermittedErr(fmt.Errorf("user id mismached"))
		return
	}

	err = s.receiptRepo.DeleteDirectionById(rDirectionId)
	if gorm.IsRecordNotFoundError(err) {
		err = nil
		return
	}

	if err != nil {
		return
	}
	return
}

func (s *Receipt) UpdateReceipt(id uint, userId uint, request UpdateReceiptRequest) (i receipt.Receipt, err error) {
	request.TrimSpaces()
	err = tools.Validator.Struct(request)
	if err != nil {
		err = tools.NewValidationErr(err)
		return
	}

	oldItem, err := s.receiptRepo.GetById(id)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}
	if err != nil {
		return
	}

	if oldItem.UserId != userId {
		err = tools.NewNotPermittedErr(fmt.Errorf("user id mismached"))
		return
	}
	i = oldItem
	i.Name = request.Name
	i.Description = request.Description
	i.Category = request.Category
	i.CookingTime = request.CookingTime
	err = s.receiptRepo.Update(&i)
	return
}

func (s *Receipt) UpdateReceiptMedia(id uint, userId uint, formFile *multipart.FileHeader, opts Options) (err error) {

	oldItem, err := s.receiptRepo.GetById(id)
	if gorm.IsRecordNotFoundError(err) {
		err = tools.NewValidationErr(fmt.Errorf("item not found"))
		return
	}
	if err != nil {
		return
	}

	if oldItem.UserId != userId {
		err = tools.NewNotPermittedErr(fmt.Errorf("user id mismached"))
		return
	}

	newMedia, err := s.mediaSvc.ProcessMedia(formFile, opts)
	if err != nil {
		return
	}

	oldItem.MediaId = &newMedia.Id
	err = s.receiptRepo.Update(&oldItem)

	return
}

func (s *Receipt) DeleteReceipt(id uint, userId uint) (err error) {
	oldItem, err := s.receiptRepo.GetById(id)
	if gorm.IsRecordNotFoundError(err) {
		err = nil
		return
	}

	if oldItem.UserId != userId {
		err = tools.NewNotPermittedErr(fmt.Errorf("user id mismatch"))
		return
	}
 	if err != nil {
		return
	}

	err = s.receiptRepo.Delete(oldItem.Id)

	return
}
