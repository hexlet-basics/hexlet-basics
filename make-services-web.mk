web-exercises-load-php:
	docker pull hexletbasics/exercises-php
	rm -rf tmp/hexletbasics/exercises-php
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-php:/out hexletbasics/exercises-php bash -c "cp -r /exercises-php/* /out"
	bundle exec rake "exercies:pull[php]"
