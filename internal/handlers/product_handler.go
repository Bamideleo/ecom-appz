package handlers

import (
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"ecom-appz/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	Repo repositories.ProductRepository
}


func NewProductHandler(repo repositories.ProductRepository) *ProductHandler{
	return  &ProductHandler{Repo: repo}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request){
	
	err := r.ParseMultipartForm(10 << 20)
	if err != nil{
		utils.JSONError(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	product := models.Product{
		Name: r.FormValue("name"),
		Description: r.FormValue("description"),
		Price: price,
		Stock: stock,
	}

	file, header, err := r.FormFile("image")

	if err == nil{
		imagePath, err := utils.SaveProductImage(file, header)
		if err != nil{
			utils.JSONError(w, "Could not save image", http.StatusInternalServerError)
			return
		}
		product.ImageURL = imagePath
	}
	
	if err := h.Repo.Create(&product); err != nil{
		utils.JSONError(w, "Could not create product", http.StatusInternalServerError)
		return
	}
	
	utils.JSONResponse(w, product, http.StatusCreated)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request){
	products, err := h.Repo.FindAll()

	if err != nil{
		utils.JSONError(w, "Could not fetch products", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, products, http.StatusOK)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request){
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	product, err := h.Repo.FindByID(id)
	if err != nil{
		utils.JSONError(w, "Product not found", http.StatusNotFound)
	}

	utils.JSONResponse(w, product, http.StatusOK)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request){
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = r.ParseMultipartForm(10 << 20)

	if err != nil{
		utils.JSONError(w, "Invalid form", http.StatusBadRequest)
		return
	}

	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	stock, _ := strconv.Atoi(r.FormValue("stock"))

	product := models.Product{
		ID: id,
		Name: r.FormValue("name"),
		Description: r.FormValue("description"),
		Price: price,
		Stock: stock,
	}
	file, header, err := r.FormFile("image")
	if err == nil{
		imagePath, err := utils.SaveProductImage(file, header)
		if err != nil{
			utils.JSONError(w, "Could not save image", http.StatusInternalServerError)
			return
		}
		product.ImageURL = imagePath
	}

	if err := h.Repo.Update(&product); err != nil{
		utils.JSONError(w, "Could not update product", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{
		"message": "Product update",
	
	}, http.StatusOK)
}


func (h *ProductHandler)Delete(w http.ResponseWriter, r *http.Request){
	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	idStr := parts[2]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	if err := h.Repo.Delete(id); err !=nil{
		utils.JSONError(w, "Could not delete product", http.StatusInternalServerError)
		return
	}
	utils.JSONResponse(w, map[string]string{
		"message":"Product deleted",
	}, http.StatusOK)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request){
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _:= strconv.Atoi(r.URL.Query().Get("limit"))

	if page <= 0{
		page =1
	}
	if limit <= 0 || limit > 100{
		limit =10
	}
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	search :=r.URL.Query().Get("search")

	products, total, err := h.Repo.List(page, limit, sort, order, search)
	if err != nil{
		utils.JSONError(w, "Could not fetch products", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"data":products,
		"meta":map[string]interface{}{
			"page":page,
			"limit":limit,
			"total":total,
			"pages": (total + limit - 1) / limit,
		},
	}
	utils.JSONResponse(w, response, http.StatusOK)
}