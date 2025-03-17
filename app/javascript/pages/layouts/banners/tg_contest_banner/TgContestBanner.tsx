import { Container } from "react-bootstrap";
import "./tg-contest-banner.scss";
import type { SharedProps } from "@/types";
import { usePage } from "@inertiajs/react";

export default function TgContestBanner() {
  const {
    props: { locale },
  } = usePage<SharedProps>();

  const localeToShowBanner = "ru";
  const currentDate = new Date();
  const endDate = new Date("2025-03-20");

  if (locale !== localeToShowBanner || currentDate > endDate) {
    return null;
  }

  return (
    <div
      className="position-sticky py-4 py-sm-3 tg-contest-banner"
      id="tg-contest-banner"
    >
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
            <a
              className="btn btn-light rounded-3 w-100 py-2 stretched-link"
              href="https://ttttt.me/hexletquizbot"
              target="_blank"
              rel="noopener noreferrer"
            >
              <span className="fs-5">Пройти тест</span>
            </a>
          </div>
        </div>
      </Container>
    </div>
  );
}
