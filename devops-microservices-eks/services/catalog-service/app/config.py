from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    # pydantic-settings reads these from env vars automatically (case-insensitive).
    # This is the FastAPI-idiomatic equivalent of auth-service's env.js —
    # one typed, validated place for config instead of os.getenv() scattered around.
    db_host: str = "localhost"
    db_port: int = 5432
    db_user: str = "postgres"
    db_password: str = "postgres"
    db_name: str = "catalogdb"

    port: int = 8000
    env: str = "development"

    # Must match the Auth service's JWT_SECRET exactly — that's how Catalog
    # trusts tokens issued by Auth without calling Auth on every request.
    jwt_secret: str
    jwt_algorithm: str = "HS256"

    @property
    def database_url(self) -> str:
        return (
            f"postgresql+psycopg2://{self.db_user}:{self.db_password}"
            f"@{self.db_host}:{self.db_port}/{self.db_name}"
        )

    class Config:
        env_file = ".env"


settings = Settings()
