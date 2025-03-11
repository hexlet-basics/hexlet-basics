import { useEffect, useState } from "react";
import { Container } from "react-bootstrap";

export default function TgContestBanner() {
  const [isVisible, setIsVisible] = useState(true);

  useEffect(() => {
    const bannerClosed = document.cookie.includes("tg_contest_banner=hide");
    if (bannerClosed) {
      setIsVisible(false);
    }
  }, []);

  const handleCloseBanner = () => {
    document.cookie = `tg_contest_banner=hide; max-age=${60 * 60 * 24}`;
    setIsVisible(false);
  };

  return (
    isVisible && (
      <div className="position-sticky py-4 py-sm-3 tg-contest-banner">
        <Container className="px-3">
          <div className="row justify-content-center justify-content-lg-end align-items-center gy-3 gy-sm-4 gy-lg-0 py-lg-1">
            <div className="col-md col-lg-7">
              <div className="fs-3 text-light lh-sm text-center text-lg-start">
                <span>Шестиклассники знают это, а ты нет😔</span>
                <span className="d-block">
                  Докажешь обратное - получишь приз!
                </span>
              </div>
            </div>
            <div className="col-sm-11 col-md-3 col-xl-2">
              <a className="btn btn-light rounded-3 w-100 py-2" href="#lalala">
                <span className="fs-5">Пройти тест</span>
              </a>
            </div>
          </div>
          <div className="position-absolute end-0 top-0 me-2 z-3">
            <button
              id="tg-contest-banner-close-btn"
              type="button"
              aria-label="Закрыть"
              className="btn-close btn-close-white"
              style={{ width: "0.5em", height: "0.5em" }}
              onClick={handleCloseBanner}
            />
          </div>
        </Container>
      </div>
    )
  );
}
