# get "*path.amp", to: redirect("/%{path}", status: 301)
get "/*path.amp", to: redirect { |params, _|
  "/#{params[:path]}"
}

# Redirects (Old routes/SEO)
get "/:locale/languages/:lang/modules/:module/lessons/:lesson", to: redirect("/%{locale}/languages/%{lang}/lessons/%{lesson}")

get "/ru/languages/pre-course-python", to: redirect("/ru/languages/python")
get "/ru/languages/pre-course-python/*", to: redirect("/ru/languages/python")

get "/ru/languages/osnovy-python", to: redirect("/ru/languages/python")
get "/ru/languages/osnovy-python/*", to: redirect("/ru/languages/python")

get "/ru/languages/layout-designer", to: redirect("/ru/languages/html")
get "/ru/languages/layout-designer/*", to: redirect("/ru/languages/html")

get "/ru/languages/pre-course-java", to: redirect("/ru/languages/java")
get "/ru/languages/pre-course-java/*", to: redirect("/ru/languages/java")

get "/ru/languages/pre-course-javascript", to: redirect("/ru/languages/javascript")
get "/ru/languages/pre-course-javascript/*", to: redirect("/ru/languages/javascript")
get "/languages/:lang/modules/:module/lessons/:lesson", to: redirect("/ru/languages/javascript")

get "/languages/:lang/modules/:module/lessons/:lesson", to: redirect("/ru/languages/%{lang}/lessons/%{lesson}")
