package handler

import (
    "strconv"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/mrmaul19/digital-library/internal/domain"
    "github.com/mrmaul19/digital-library/pkg/response"
    "github.com/mrmaul19/digital-library/pkg/upload"
)

type BookHandler struct {
    bookUsecase domain.BookUsecase
    uploader    *upload.Uploader
}

func NewBookHandler(bookUsecase domain.BookUsecase, uploader *upload.Uploader) *BookHandler {
    return &BookHandler{bookUsecase: bookUsecase, uploader: uploader}
}

func (h *BookHandler) Create(c *fiber.Ctx) error {
    var req domain.CreateBookRequest

    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "format request tidak valid")
    }

    if req.Title == "" || req.Author == "" || req.Genre == "" {
        return response.BadRequest(c, "title, author, dan genre wajib diisi")
    }

    // Upload cover dan file buku
    coverURL, err := h.uploader.UploadCover(c, "cover")
    if err != nil {
        return response.BadRequest(c, err.Error())
    }

    fileURL, err := h.uploader.UploadBook(c, "file")
    if err != nil {
        return response.BadRequest(c, err.Error())
    }

    // Ambil userID dari JWT middleware
    userID := c.Locals("userID").(uuid.UUID)

    book, err := h.bookUsecase.Create(&req, coverURL, fileURL, userID)
    if err != nil {
        return response.InternalError(c, err.Error())
    }

    return response.Created(c, "buku berhasil ditambahkan", book)
}

func (h *BookHandler) GetAll(c *fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    year, _ := strconv.Atoi(c.Query("year", "0"))

    filter := domain.BookFilter{
        Search: c.Query("search"),
        Genre:  c.Query("genre"),
        Author: c.Query("author"),
        Year:   year,
        Page:   page,
        Limit:  limit,
    }

    result, err := h.bookUsecase.GetAll(filter)
    if err != nil {
        return response.InternalError(c, err.Error())
    }

    return response.Paginated(c, "berhasil mengambil data buku", result.Books, response.Meta{
        Total:      result.Total,
        Page:       result.Page,
        Limit:      result.Limit,
        TotalPages: result.TotalPages,
    })
}

func (h *BookHandler) GetByID(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.BadRequest(c, "ID tidak valid")
    }

    book, err := h.bookUsecase.GetByID(id)
    if err != nil {
        return response.NotFound(c, err.Error())
    }

    return response.OK(c, "berhasil mengambil data buku", book)
}

func (h *BookHandler) Update(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.BadRequest(c, "ID tidak valid")
    }

    var req domain.UpdateBookRequest
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "format request tidak valid")
    }

    book, err := h.bookUsecase.Update(id, &req)
    if err != nil {
        return response.NotFound(c, err.Error())
    }

    return response.OK(c, "buku berhasil diupdate", book)
}

func (h *BookHandler) Delete(c *fiber.Ctx) error {
    id, err := uuid.Parse(c.Params("id"))
    if err != nil {
        return response.BadRequest(c, "ID tidak valid")
    }

    if err := h.bookUsecase.Delete(id); err != nil {
        return response.NotFound(c, err.Error())
    }

    return response.OK(c, "buku berhasil dihapus", nil)
}