from pydantic import BaseModel, Field


class ProductIn(BaseModel):
    name: str = Field(min_length=1, max_length=255)
    description: str | None = Field(default=None, max_length=1000)
    price: float = Field(gt=0, description="must be greater than 0")
    stock: int = Field(default=0, ge=0)


class ProductUpdate(BaseModel):
    # All optional — PATCH-style partial update
    name: str | None = Field(default=None, min_length=1, max_length=255)
    description: str | None = Field(default=None, max_length=1000)
    price: float | None = Field(default=None, gt=0)
    stock: int | None = Field(default=None, ge=0)


class ProductOut(ProductIn):
    id: int

    class Config:
        from_attributes = True


class ErrorResponse(BaseModel):
    error: str
    message: str
