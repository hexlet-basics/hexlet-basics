import { Accordion, Alert, Card, Col, Container, Row } from "react-bootstrap";

import learningEnVideo from "@/images/course-landing-page/learning_en.mp4";
import learningRuVideo from "@/images/course-landing-page/learning_ru.mp4";
import waitingClock from "@/images/waiting_clock.png";
import { useTranslation } from "react-i18next";

import XssContent from "@/components/XssContent";
import ApplicationLayout from "@/pages/layouts/ApplicationLayout";
import * as Routes from "@/routes.js";
import type { BreadcrumbItem, Language, SharedProps } from "@/types";
import type {
  LanguageCategory,
  LanguageLandingPage,
  LanguageLandingPageQnaItem,
  LanguageLesson,
  LanguageMember,
  LanguageModule,
} from "@/types/serializers";
import { Head, Link, usePage } from "@inertiajs/react";
import type { Product, WithContext } from "schema-dts";

type Props = {
  course: Language;
  courseMember?: LanguageMember;
  courseCategory: LanguageCategory;
  courseLandingPage: LanguageLandingPage;
  courseLandingPageQnaItems: LanguageLandingPageQnaItem[];
  firstLesson: LanguageLesson;
  nextLesson?: LanguageLesson;
  // user: User;
  courseModules: LanguageModule[];
  lessonsByModuleId: {
    [moduleId: number]: LanguageLesson[];
  };
};

export default function Show({
  firstLesson,
  nextLesson,
  courseLandingPage,
  courseLandingPageQnaItems,
  courseMember,
  courseCategory,
  courseModules,
  lessonsByModuleId,
  course,
}: Props) {
  const { t } = useTranslation();
  const { auth, locale } = usePage<SharedProps>().props;

  const product: WithContext<Product> = {
    "@context": "https://schema.org",
    "@type": "Product",
    description: courseLandingPage.description,
    image: course.cover_list_variant,
    offers: {
      "@type": "Offer",
      price: "0",
      priceCurrency: "RUB",
      availability: "https://schema.org/InStock",
    },
    aggregateRating: {
      "@type": "AggregateRating",
      ratingValue: course.rating_value,
      ratingCount: course.rating_count,
    },
    name: courseLandingPage.header,
    // TODO: add review
  };

  const breadcrumbItems: BreadcrumbItem[] = [
    {
      name: courseCategory.name,
      url: Routes.language_category_url(courseCategory.slug!),
    },
    {
      name: courseLandingPage.header,
      url: Routes.language_url(courseLandingPage.slug),
    },
  ];

  return (
    <ApplicationLayout items={breadcrumbItems}>
      <Head>
        <script type="application/ld+json">{JSON.stringify(product)}</script>
      </Head>
      <Container className="pt-4 pt-lg-5">
        {courseMember?.state === "finished" && (
          <Alert variant="success">
            <XssContent>{t("languages.show.completed_html")}</XssContent>
          </Alert>
        )}
        <Row className="justify-content-center py-2 mb-5">
          <Col className="col-lg-7 text-center">
            <div className="fs-5 fw-medium text-primary text-opacity-75 lh-sm mb-2">
              {t("languages.show.free_course")}
            </div>
            <h1 className="display-5 fw-bolder mb-3">
              {courseLandingPage.header}
            </h1>
            <div className="fs-5 text-body-secondary mb-5">
              {courseLandingPage.description}
            </div>
            <Row className="row-cols-1 row-cols-sm-2 gy-3">
              {!courseMember && (
                <Col className="col-lg-4">
                  <Link
                    className="btn btn-primary w-100"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      firstLesson.slug,
                    )}
                  >
                    <span>{t("languages.show.start")}</span>
                  </Link>
                </Col>
              )}
              {courseMember && nextLesson && (
                <Col className="col-lg-4">
                  <Link
                    className="btn btn-outline-primary w-100"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      nextLesson.slug,
                    )}
                  >
                    <span>{t("languages.show.continue")}</span>
                  </Link>
                </Col>
              )}
              {auth.user.guest && (
                <Col className="col-lg-5">
                  <Link className="btn fw-medium" href={Routes.new_user_path()}>
                    <span className="me-2">
                      {t("languages.show.registration")}
                    </span>
                    <i className="bi bi-arrow-right" />
                  </Link>
                </Col>
              )}
            </Row>
          </Col>
        </Row>
        <Row className="mb-lg-5 py-4 py-md-5">
          <Col className="col-lg-9">
            <h2 className="display-5 fw-medium lh-1 mb-4">
              {courseLandingPage.used_in_header}
            </h2>
            <p className="text-body-secondary pe-lg-2">
              {courseLandingPage.used_in_description}
            </p>
          </Col>
        </Row>
        <Row className="row-cols-1 row-cols-md-2 gx-lg-5 gy-5 mb-lg-4 py-3 py-md-5 justify-content-between">
          <Col>
            {courseLandingPage.outcomes_image && (
              <img
                src={courseLandingPage.outcomes_image}
                width="100%"
                height="auto"
                alt={t("languages.show.learning_preview")}
                className="rounded-5 shadow-lg"
              />
            )}
          </Col>
          <Col>
            <div className="display-5 fw-semibold lh-1 mb-4">
              {courseLandingPage.outcomes_header}
            </div>
            <p className="text-body-secondary pe-lg-5">
              {courseLandingPage.outcomes_description}
            </p>
          </Col>
        </Row>
        <Row className="mb-lg-5 py-5">
          <Col className="col-lg-10">
            <div className="display-5 fw-semibold lh-1 mb-4">
              {t("languages.show.learning_program")}
            </div>
            <Accordion defaultActiveKey="0" className="hexlet-basics-accordion">
              {courseModules.map((m, index) => (
                <Accordion.Item
                  eventKey={index.toString()}
                  key={m.id}
                  className="rounded-0 border-0 border-bottom border-secondary-subtle py-3 py-md-4"
                >
                  <Accordion.Header as="h3">
                    {index + 1}. {m.name!}
                  </Accordion.Header>
                  <Accordion.Body className="px-0 pb-0">
                    <ul className="list-unstyled">
                      {(lessonsByModuleId[m.id] ?? []).map((l) => (
                        <li key={l.id}>
                          <Link
                            className="text-decoration-none text-body-secondary"
                            href={Routes.language_lesson_path(
                              course.slug!,
                              l.slug!,
                            )}
                          >
                            <span>{l.name}</span>
                          </Link>
                        </li>
                      ))}
                    </ul>
                  </Accordion.Body>
                </Accordion.Item>
              ))}
            </Accordion>
          </Col>
        </Row>
        <div className="display-5 fw-semibold lh-1 mb-4 mb-lg-0">
          {t("languages.show.about_learning")}
        </div>
        <Row className="row-cols-1 row-cols-lg-2 mb-lg-5 pb-4 py-md-5 gy-4">
          <Col>
            <div className="me-lg-5 pe-lg-4">
              <div className="d-flex mb-3">
                <i className="bi bi-cloud-arrow-up-fill me-3 text-primary" />
                <XssContent>
                  {t("languages.show.without_registration_html")}
                </XssContent>
              </div>
              <div className="d-flex mb-3">
                <i className="bi bi-laptop-fill me-3 text-primary" />
                <XssContent>
                  {t("languages.show.learning_conveniently_html")}
                </XssContent>
              </div>
              <div className="d-flex">
                <i className="bi bi-lock-fill me-3 text-primary" />
                <XssContent>
                  {t("languages.show.real_life_challenges_html")}
                </XssContent>
              </div>
            </div>
          </Col>
          <Col>
            <div className="bg-primary rounded-5 overflow-hidden py-lg-3">
              <video
                className="w-100 rounded-4 hexlet-basics-learning-video"
                src={locale === "en" ? learningEnVideo : learningRuVideo}
                autoPlay
                loop
                muted
                playsInline
              />
            </div>
          </Col>
        </Row>
        <div className="py-4 mb-lg-5">
          <Card className="px-3 px-sm-4 px-lg-5 py-5">
            <Row className="gy-4 justify-content-between align-items-center py-lg-5">
              <Col className="col-12 col-sm-11 col-lg-9 col-xxl-7">
                <div className="display-5 fw-semibold lh-1">
                  Готовы к инновациям в обучении? Начните обучение с AI прямо
                  сейчас
                </div>
              </Col>
              <Col className="col-auto">
                {!courseMember && (
                  <Link
                    className="btn btn-primary"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      firstLesson.slug,
                    )}
                  >
                    <span>{t("languages.show.start")}</span>
                  </Link>
                )}
                {courseMember && nextLesson && (
                  <Link
                    className="btn btn-primary"
                    href={Routes.language_lesson_path(
                      course.slug!,
                      nextLesson.slug,
                    )}
                  >
                    <span>{t("languages.show.continue")}</span>
                  </Link>
                )}
              </Col>
            </Row>
          </Card>
        </div>
        <div className="py-4 py-lg-5">
          <div className="display-5 fw-semibold lh-1 mb-5">Отзывы</div>
          <Row className="row-cols-1 row-cols-lg-2 row-cols-xl-3 gy-4">
            <Col>
              <div className="d-flex flex-column h-100 rounded-3 bg-body-tertiary p-4">
                <div className="mb-3">
                  Простота, легкость, доходчивость. Я, человек который никогда
                  не занимался и не изучал программирование (не считая уроков
                  информатики), очень легко вошел в суть базы Python. Да, я
                  только начал, даже не закончил программу, но на моей памяти
                  сохранилось все то, когда только начинал. Мне понравилось, как
                  я (а-ля студент) за счет простых и часто интересных примеров
                  изучаю темы. Я понимаю, что здесь не смогу полностью понять и
                  стать гуру в Puthon-е, но для старта и понимания "интересна ли
                  тебе эта профессия?" - я считаю этот сайт отличным примером
                </div>
                <div className="mt-auto">
                  <div className="fw-bold">Timur</div>
                  <div>Сcылка на курс</div>
                </div>
              </div>
            </Col>
            <Col>
              <div className="d-flex flex-column h-100 rounded-3 bg-body-tertiary p-4">
                <div className="mb-3">
                  Удобство интерфейса Платформа имеет простой и интуитивно
                  понятный интерфейс, что позволяет легко ориентироваться среди
                  курсов и материалов. Качество контента Курсы содержат
                  актуальные и полезные материалы. Они прекрасно
                  структурированы, что облегчает усвоение информации
                </div>
                <div className="mt-auto">
                  <div className="fw-bold">Цапкова Любовь Евгеньевна</div>
                  <div>Сcылка на курс</div>
                </div>
              </div>
            </Col>
            <Col>
              <div className="d-flex flex-column h-100 rounded-3 bg-body-tertiary p-4">
                <div className="mb-3">
                  Очень понравился курс. Задачи по возрастанию сложности,
                  понятное объяснение тем. Удобная песочница и в конце можно
                  посмотреть вариант решения. Он не всегда совпадает с моим.
                  Удобно, что можно изучить тесты и это помогает в решении. так
                  же помогает, что можно посмотреть почему не прошли тесты и
                  исправить свое решение
                </div>
                <div className="mt-auto">
                  <div className="fw-bold">Анна Задирака</div>
                  <div>Сcылка на курс</div>
                </div>
              </div>
            </Col>
            <Col>
              <div className="d-flex flex-column h-100 rounded-3 bg-body-tertiary p-4">
                <div className="mb-3">
                  Курс был организован очень логично и последовательно. Он
                  начинался с основ, что идеально подходит для новичков, и
                  постепенно переходил к более сложным темам. Каждый модуль был
                  четко структурирован, что позволяло легко следовать за
                  материалом и не терять нить
                </div>
                <div className="mt-auto">
                  <div className="fw-bold">Дмитрий Кондрашов</div>
                  <div>Сcылка на курс</div>
                </div>
              </div>
            </Col>
            <Col>
              <div className="d-flex flex-column h-100 rounded-3 bg-body-tertiary p-4">
                <div className="mb-3">
                  Очень грамотно и мягко вводят в курс HTML (Проходил конкретно
                  его). Дают очень хорошую базу, подкрепляют задачами в конце
                  каждого урока. Полностью стоит своего времени. Я как человек с
                  нулевым бэкграундом в верстке, остался полностью доволен
                  данным курсом
                </div>
                <div className="mt-auto">
                  <div className="fw-bold">Игорь</div>
                  <div>Сcылка на курс</div>
                </div>
              </div>
            </Col>
            <Col>
              <div className="d-flex flex-column h-100 rounded-3 bg-body-tertiary p-4">
                <div className="mb-3">
                  Сначала относился скептически, потому что базовых курсов для
                  новичков хоть пруд-пруди. Но когда начал проходить обратил
                  внимание на хорошие примеры использования кода и задачки после
                  каждого урока, заставляющие хоть немножко задуматься. Хоть
                  базовые знания уже и имеются, но какие-то моменты узнал
                  впервые. Разборался с тернальными операторами, свойствами и
                  методами и с интерполяций, что вызывало у меня какие-то
                  затруднения
                </div>
                <div className="mt-auto">
                  <div className="fw-bold">Виталий</div>
                  <div>Сcылка на курс</div>
                </div>
              </div>
            </Col>
          </Row>
        </div>
        <div className="py-4 py-lg-5">
          <div className="bg-dark text-light p-4 rounded-3">
            <Row className="flex-column flex-lg-row">
              <Col className="col-lg-4 col-xl-5 col-xxl-6">
                <img
                  src={waitingClock}
                  width="100%"
                  height="auto"
                  alt={t("languages.show.learning_preview")}
                />
              </Col>
              <Col className="d-flex flex-column justify-content-center py-4 py-md-5">
                <div className="display-5 fw-semibold lh-1 mb-4">
                  Больше чем Поддержка
                </div>
                <div className="pe-lg-5">
                  <div className="mb-5">
                    Мы знаем, как непросто начинать в IT, поэтому создали
                    сообщество разработчиков, где вам всегда готовы помочь.
                    Здесь можно задавать вопросы, получать поддержку, общаться с
                    опытными специалистами и быстрее влиться в профессию
                  </div>
                  <Link
                    className="btn btn-outline-secondary"
                    href="https://ttttt.me/HexletLearningBot"
                  >
                    <span>Присоединиться</span>
                  </Link>
                </div>
              </Col>
            </Row>
          </div>
        </div>
        {courseLandingPageQnaItems.length > 0 && (
          <div className="py-4 py-lg-5">
            <div className="display-5 fw-semibold lh-1 mb-5">
              {t("languages.show.sort_questions")}
            </div>
            <Row className="gy-4 gy-md-5">
              {courseLandingPageQnaItems.map((item) => (
                <Col key={item.id} className="col-12 col-md-6">
                  <div className="fs-5 fw-medium mb-3 pe-lg-5">
                    {item.question}
                  </div>
                  <p className="pe-lg-5">{item.answer}</p>
                </Col>
              ))}
            </Row>
          </div>
        )}
      </Container>
    </ApplicationLayout>
  );
}
