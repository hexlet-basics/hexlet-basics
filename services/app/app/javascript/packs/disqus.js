import gon from 'gon';

window.disqus_config = () => {
  this.page.url = window.location.href;
  this.page.identifier = `lesson-${gon.lesson.id}`;
};

window.disqus_shortname = 'hexlet-basics';
window.disqus_identifier = `lesson-${gon.lesson.id}`;
window.disqus_url = window.location.href;

(() => {
  const d = document;
  const s = d.createElement('script');
  s.src = 'https://hexlet-basics.disqus.com/embed.js';
  s.setAttribute('data-timestamp', +new Date());
  (d.head || d.body).appendChild(s);
})();

window.reset = (newIdentifier, newUrl, newTitle, newLanguage) => {
  window.DISQUS.reset({
    reload: true,
    config() {
      this.page.identifier = newIdentifier;
      this.page.url = newUrl;
      this.page.title = newTitle;
      this.language = newLanguage;
    },
  });
};
