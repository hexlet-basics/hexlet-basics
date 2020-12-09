import gon from 'gon';

window.disqus_config = function () {
  this.page.url = window.location.href;  // Replace PAGE_URL with your page's canonical URL variable
  this.page.identifier = `lesson-${gon.lesson.id}`; // Replace PAGE_IDENTIFIER with your page's unique identifier variable
};

window.disqus_shortname = 'hexlet-basics';
window.disqus_identifier = `lesson-${gon.lesson.id}`;
window.disqus_url = window.location.href;

(function () { // DON'T EDIT BELOW THIS LINE
  var d = document, s = d.createElement('script');
  s.src = 'https://hexlet-basics.disqus.com/embed.js';
  s.setAttribute('data-timestamp', +new Date());
  (d.head || d.body).appendChild(s);
})();

window.reset = function (newIdentifier, newUrl, newTitle, newLanguage) {
  window.DISQUS.reset({
    reload: true,
    config: function () {
      this.page.identifier = newIdentifier;
      this.page.url = newUrl;
      this.page.title = newTitle;
      this.language = newLanguage;
    }
  });
};
