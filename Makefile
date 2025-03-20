IMAGE_NAME = iic3413-evaluator

build:
	@docker build -t $(IMAGE_NAME) .
run:
	@docker run --rm \
		-v ./io:/app/io \
		-v ./config:/app/config \
		$(IMAGE_NAME)
