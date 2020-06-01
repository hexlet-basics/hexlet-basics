web-exercises-load-php:
	docker pull hexletbasics/exercises-php
	rm -rf tmp/hexletbasics/exercises-php
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-php:/out hexletbasics/exercises-php bash -c "cp -r /exercises-php/* /out"
	bundle exec rake "exercises:load[php]"

web-exercises-load-python:
	docker pull hexletbasics/exercises-python
	rm -rf tmp/hexletbasics/exercises-python
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-python:/out hexletbasics/exercises-python bash -c "cp -r /exercises-python/* /out"
	bundle exec rake "exercises:load[python]"

web-exercises-load-ruby:
	docker pull hexletbasics/exercises-ruby
	rm -rf tmp/hexletbasics/exercises-ruby
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-ruby:/out hexletbasics/exercises-ruby bash -c "cp -r /exercises-ruby/* /out"
	bundle exec rake "exercises:load[ruby]"

web-exercises-load-racket:
	docker pull hexletbasics/exercises-racket
	rm -rf tmp/hexletbasics/exercises-racket
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-racket:/out hexletbasics/exercises-racket bash -c "cp -r /exercises-racket/* /out"
	bundle exec rake "exercises:load[racket]"

web-exercises-load-javascript:
	docker pull hexletbasics/exercises-javascript
	rm -rf tmp/hexletbasics/exercises-javascript
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-javascript:/out hexletbasics/exercises-javascript bash -c "cp -r /exercises-javascript/* /out"
	bundle exec rake "exercises:load[javascript]"

web-exercises-load-java:
	docker pull hexletbasics/exercises-java
	rm -rf tmp/hexletbasics/exercises-java
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-java:/out hexletbasics/exercises-java bash -c "cp -r /exercises-java/* /out"
	bundle exec rake "exercises:load[java]"

web-exercises-load-html:
	docker pull hexletbasics/exercises-html
	rm -rf tmp/hexletbasics/exercises-html
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-html:/out hexletbasics/exercises-html bash -c "cp -r /exercises-html/* /out"
	bundle exec rake "exercises:load[html]"

web-exercises-load-css:
	docker pull hexletbasics/exercises-css
	rm -rf tmp/hexletbasics/exercises-css
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-css:/out hexletbasics/exercises-css bash -c "cp -r /exercises-css/* /out"
	bundle exec rake "exercises:load[css]"

web-exercises-load-elixir:
	docker pull hexletbasics/exercises-elixir
	rm -rf tmp/hexletbasics/exercises-elixir
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-elixir:/out hexletbasics/exercises-elixir bash -c "cp -r /exercises-elixir/* /out"
	bundle exec rake "exercises:load[elixir]"

web-exercises-load-go:
	docker pull hexletbasics/exercises-go
	rm -rf tmp/hexletbasics/exercises-go
	docker run --rm -v $(CURDIR)/tmp/hexletbasics/exercises-go:/out hexletbasics/exercises-go bash -c "cp -r /exercises-go/* /out"
	bundle exec rake "exercises:load[go]"
