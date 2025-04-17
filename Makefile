IMAGE_NAME = iic3413-evaluator

build:
	@docker build -t $(IMAGE_NAME) .
run:
	@cp -r io io-cpy
	@chmod go-wrx io-cpy
	@docker run --rm -v ./io-cpy:/app/io-cpy \
		$(IMAGE_NAME) -n=${LAB_NAME} -m=${MODE} || :
	@cp io-cpy/results/* io/results/
	@chmod go+r io/results/*
	@rm -rf io-cpy
