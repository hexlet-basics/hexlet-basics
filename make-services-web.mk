web-exercises-load-javascript:
	docker pull hexletbasics/exercises-javascript
	rm -rf tmp/hexletbasics/exercises-javascript
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-javascript:/out hexletbasics/exercises-javascript bash -c "cp -r /exercises-javascript/* /out"
	# docker-compose run --rm web bundle exec bin/rails runner exercises.rb javascript
