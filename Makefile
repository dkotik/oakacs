default:
	cd v1 && for s in $$(go list ./...); do if ! go test -failfast -v -p 1 $$s; then break; fi; done
