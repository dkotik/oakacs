default:
	cd oakrbac && go test
	@# cd v1/oakwords && go test -v ./...
	@# cd v1 && for s in $$(go list ./...); do if ! go test -failfast -v -p 1 -timeout 60s $$s; then break; fi; done
