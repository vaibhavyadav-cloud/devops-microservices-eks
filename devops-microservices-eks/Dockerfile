
FROM public.ecr.aws/lambda/python:3.11

WORKDIR /var/task
COPY requirements.txt .

# Prefer wheels and refuse source builds for PyMuPDF
ENV PIP_ONLY_BINARY=:all: \
    PIP_PREFER_BINARY=1

RUN python -m pip install --upgrade pip setuptools wheel && \
    pip install --no-cache-dir -r requirements.txt

COPY . .
CMD ["lambda_handler.lambda_handler"]
