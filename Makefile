run:
	go run . ./examples/$(filter-out $@,$(MAKECMDGOALS)).ser

%:
	@:
