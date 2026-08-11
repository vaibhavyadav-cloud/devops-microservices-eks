from fastapi import APIRouter, Depends, Query
from sqlalchemy.orm import Session

from ..database import get_db, Product
from ..schemas import ProductIn, ProductOut, ProductUpdate
from ..exceptions import ProductNotFoundError
from ..logger import logger
from ..middleware.auth import get_current_user

router = APIRouter(prefix="/products", tags=["products"])


@router.post("", response_model=ProductOut, status_code=201)
def create_product(
    product: ProductIn,
    db: Session = Depends(get_db),
    user: dict = Depends(get_current_user),  # protected — must be logged in
):
    db_product = Product(**product.model_dump())
    db.add(db_product)
    db.commit()
    db.refresh(db_product)
    logger.info(f"product_created id={db_product.id} by_user_id={user['userId']}")
    return db_product


@router.get("", response_model=list[ProductOut])
def list_products(
    db: Session = Depends(get_db),
    limit: int = Query(default=20, ge=1, le=100),
    offset: int = Query(default=0, ge=0),
):
    # Pagination — without this, a products table with 100k rows would
    # blow the response payload and the DB round-trip time.
    return db.query(Product).offset(offset).limit(limit).all()


@router.get("/{product_id}", response_model=ProductOut)
def get_product(product_id: int, db: Session = Depends(get_db)):
    product = db.query(Product).filter(Product.id == product_id).first()
    if not product:
        raise ProductNotFoundError(product_id)
    return product


@router.patch("/{product_id}", response_model=ProductOut)
def update_product(
    product_id: int,
    payload: ProductUpdate,
    db: Session = Depends(get_db),
    user: dict = Depends(get_current_user),  # protected — must be logged in
):
    product = db.query(Product).filter(Product.id == product_id).first()
    if not product:
        raise ProductNotFoundError(product_id)

    updates = payload.model_dump(exclude_unset=True)
    for field, value in updates.items():
        setattr(product, field, value)

    db.commit()
    db.refresh(product)
    logger.info(f"product_updated id={product.id} by_user_id={user['userId']}")
    return product


@router.delete("/{product_id}", status_code=204)
def delete_product(
    product_id: int,
    db: Session = Depends(get_db),
    user: dict = Depends(get_current_user),  # protected — must be logged in
):
    product = db.query(Product).filter(Product.id == product_id).first()
    if not product:
        raise ProductNotFoundError(product_id)
    db.delete(product)
    db.commit()
    logger.info(f"product_deleted id={product_id} by_user_id={user['userId']}")
