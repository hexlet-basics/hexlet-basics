// @ts-check

const bannerCloseBtn = document.getElementById('competition-banner-close-btn');
const banner = document.getElementById('competition-banner');

bannerCloseBtn.addEventListener('click', (event) => {
  event.preventDefault();
  document.cookie = `competition_banner=hide; max-age=${60 * 60 * 24}`;
  banner.style.transition = 'opacity 0.3s';
  banner.style.opacity = '0';

  setTimeout(() => {
    banner.style.display = 'none';
  }, 300);
});
