import * as Routes from '@/routes.js';

export default function CookiePolicyEn() {
  return (
    <>
      <p className="fw-medium">Last updated March 28, 2025</p>

      <p>
        This Cookie Policy explains how TOO "Hexlet" (
        <strong>"Company," "we," "us,"</strong> and
        <strong>"our"</strong>) uses cookies and similar technologies to
        recognize you when you visit our website at{' '}
        <a href={Routes.root_path()} target="_blank" rel="noreferrer">
          https://code-basics.com
        </a>{' '}
        (<strong>"Website"</strong>). It explains what these technologies are
        and why we use them, as well as your rights to control our use of them.
      </p>

      <p>
        In some cases we may use cookies to collect personal information, or
        that becomes personal information if we combine it with other
        information.
      </p>

      <strong>What are cookies?</strong>

      <p>
        Cookies are small data files that are placed on your computer or mobile
        device when you visit a website. Cookies are widely used by website
        owners in order to make their websites work, or to work more
        efficiently, as well as to provide reporting information.
      </p>

      <p>
        Cookies set by the website owner (in this case, TOO "Hexlet") are called
        "first-party cookies." Cookies set by parties other than the website
        owner are called "third-party cookies." Third-party cookies enable
        third-party features or functionality to be provided on or through the
        website (e.g., advertising, interactive content, and analytics). The
        parties that set these third-party cookies can recognize your computer
        both when it visits the website in question and also when it visits
        certain other websites.
      </p>

      <strong>Why do we use cookies?</strong>

      <p>
        We use first- and third-party cookies for several reasons. Some cookies
        are required for technical reasons in order for our Website to operate,
        and we refer to these as "essential" or "strictly necessary" cookies.
        Other cookies also enable us to track and target the interests of our
        users to enhance the experience on our Online Properties. Third parties
        serve cookies through our Website for advertising, analytics, and other
        purposes. This is described in more detail below
      </p>

      <strong>How can I control cookies?</strong>

      <p>
        You have the right to decide whether to accept or reject cookies. You
        can exercise your cookie rights by setting your preferences in the
        Cookie Consent Manager. The Cookie Consent Manager allows you to select
        which categories of cookies you accept or reject. Essential cookies
        cannot be rejected as they are strictly necessary to provide you with
        services.
      </p>

      <p>
        The Cookie Consent Manager can be found in the notification banner and
        on our Website. If you choose to reject cookies, you may still use our
        Website though your access to some functionality and areas of our
        Website may be restricted. You may also set or amend your web browser
        controls to accept or refuse cookies.
      </p>

      <p>
        The specific types of first- and third-party cookies served through our
        Website and the purposes they perform are described in the table below
        (please note that the specific cookies served may vary depending on the
        specific Online Properties you visit):
      </p>

      <strong className="text-decoration-underline">
        Performance and functionality cookies:
      </strong>

      <p>
        These cookies are used to enhance the performance and functionality of
        our Website but are non-essential to their use. However, without these
        cookies, certain functionality (like videos) may become unavailable.
      </p>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>XSRF-TOKEN</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                This cookie is written to help with site security in preventing
                Cross-Site Request Forgery attacks.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                Advertiser's website domain{' '}
                <a
                  href={Routes.page_path('cookie_policy')}
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>session</td>
            </tr>
          </tbody>
        </table>
      </div>

      <strong className="text-decoration-underline">
        Analytics and customization cookies:
      </strong>

      <p>
        These cookies collect information that is used either in aggregate form
        to help us understand how our Website is being used or how effective our
        marketing campaigns are, or to help us customize our Website for you.
      </p>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>s7</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Gather data regarding site usage and user behavior on the
                website.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td> .yandex.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>Adobe Analytics</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>11 months 30 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>_ym_uid</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Used by Yandex Metrica as a unique user ID to help track a user
                in a website
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                Metrica{' '}
                <a
                  href="https://yandex.com/legal/confidential/?lang=en"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>11 months 30 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>_ga</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Records a particular ID used to come up with data about website
                usage by the user
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                Google Analytics{' '}
                <a
                  href="https://business.safety.google/privacy/"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 year 1 month 4 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>s7</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Gather data regarding site usage and user behavior on the
                website.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>Adobe Analytics</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_local_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td> persistent</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>_ym_uid</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Used by Yandex Metrica as a unique user ID to help track a user
                in a website
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                Metrica{' '}
                <a
                  href="https://yandex.com/legal/confidential/?lang=en"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_local_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>persistent</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>s7</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Gather data regarding site usage and user behavior on the
                website.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>Adobe Analytics</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>59 minutes</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>s7</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Gather data regarding site usage and user behavior on the
                website.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.yandex.ru</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td> Adobe Analytics</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>11 months 30 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>s7</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Gather data regarding site usage and user behavior on the
                website.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.mc.yandex.ru</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>Adobe Analytics</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>10 minutes</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>s7</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Gather data regarding site usage and user behavior on the
                website.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>prism.app-us1.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>Adobe Analytics</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>30 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>tmr_detect</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                This is a Vk plugin cookie used for profiling users. It expires
                at the end of a browsing session.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                VK{' '}
                <a
                  href="https://vk.com/privacy"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>23 hours 59 minutes</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>1</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>mc.yandex.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>pixel_tracker</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>session</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>s7</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Gather data regarding site usage and user behavior on the
                website
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>mc.yandex.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>Adobe Analytics</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>pixel_tracker</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>session</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>_ym_isad</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Used by Yandex Metrica to determine if a visitor has ad blockers
                installed in their browsers
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                Metrica{' '}
                <a
                  href="https://yandex.com/legal/confidential/?lang=en"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>19 hours 59 minutes</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>_ym_d</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Used by Yandex Metrica to determine the date of the user's first
                site session.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                Metrica{' '}
                <a
                  href="https://yandex.com/legal/confidential/?lang=en"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>11 months 30 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>54648685</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>mc.yandex.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>pixel_tracker</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>session</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>s7</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Gather data regarding site usage and user behavior on the
                website.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.mc.yandex.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>Adobe Analytics</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 day</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>i</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Used by Yandex Metrica for identifying site users. This cookie
                exists for a period of 1 year maximum.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>yandex.ru</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                Yandex Metrica{' '}
                <a
                  href="https://yandex.com/legal/confidential/?lang=en"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 year 11 months 29 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>_ga_#</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Used to distinguish individual users by means of designation of
                a randomly generated number as client identifier, which allows
                calculation of visits and sessions
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                Google Analytics{' '}
                <a
                  href="https://business.safety.google/privacy/"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 year 1 month 4 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>s7</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Gather data regarding site usage and user behavior on the
                website.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.vk.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>Adobe Analytics</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td> 11 months 30 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <strong className="text-decoration-underline">
        Advertising cookies:
      </strong>

      <p>
        These cookies are used to make advertising messages more relevant to
        you. They perform functions like preventing the same ad from
        continuously reappearing, ensuring that ads are properly displayed for
        advertisers, and in some cases selecting advertisements that are based
        on your interests.
      </p>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>IDE</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Used to measure the conversion rate of ads presented to the
                user. Expires in 1.5 years.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.doubleclick.net</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                DoubleClick{' '}
                <a
                  href="https://business.safety.google/privacy/"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 year 11 months 29 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>test_cookie</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                A session cookie used to check if the user’s browser supports
                cookies.
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.doubleclick.net</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                DoubleClick{' '}
                <a
                  href="https://business.safety.google/privacy/"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>15 minutes</td>
            </tr>
          </tbody>
        </table>
      </div>

      <strong className="text-decoration-underline">
        Social networking cookies:
      </strong>

      <p>
        These cookies are used to enable you to share pages and content that you
        find interesting on our Website through third-party social networking
        and other websites. These cookies may also be used for advertising
        purposes.
      </p>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>rtrg</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Is a social tracking cookie used by VKontakte. This cookie
                expires in 368 days
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>vk.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                VK{' '}
                <a
                  href="https://vk.com/privacy"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>pixel_tracker</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>session</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>remixlang</td>
            </tr>
            <tr>
              <td>Purpose:</td>
              <td>
                Is a social tracking cookie used by VKontakte. This cookie
                expires in 368 days
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.vk.com</td>
            </tr>
            <tr>
              <td>Service:</td>
              <td>
                VK{' '}
                <a
                  href="https://vk.com/privacy"
                  target="_blank"
                  rel="noreferrer"
                >
                  View Service Privacy Policy
                </a>
              </td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 year 4 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <strong className="text-decoration-underline">
        Unclassified cookies:
      </strong>

      <p>
        These are cookies that have not yet been categorized. We are in the
        process of classifying these cookies with the help of their providers.
      </p>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>__ym_tab_guid</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_session_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>session</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>tmr_lvidTS</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_local_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>persistent</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>__anon_id</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_local_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>persistent</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>ymex</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.yandex.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>11 months 30 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>tmr_reqNum</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_local_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>persistent</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>_ym_wv2rf:54648685:0</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_local_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>persistent</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>tmr_trgfpid</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_local_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>persistent</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>yandexuid</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.yandex.ru</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 year 11 months 29 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>bh</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.yandex.ru</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 year 1 month 4 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>_ym54648685:0_reqNum</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_local_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>persistent</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>bh</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td> .yandex.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 year 1 month 4 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>tmr_lvid</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>10 months 28 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>
                ph_phc_y7Jf3JmTEdmXtqZuHoKFXS2R8ypILrhna4VV1IpmrGU_window_id
              </td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_session_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>session</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>tmr_lvid</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>html_local_storage</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td> persistent</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>tmr_lvidTS</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.code-basics.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>http_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>10 months 28 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>yp</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.yandex.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>9 years 11 months 28 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>yandexuid</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.yandex.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>1 year 11 months 29 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div className="table-responsive mb-3">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <td>Name:</td>
              <td>receive-cookie-deprecation</td>
            </tr>
            <tr>
              <td>Provider:</td>
              <td>.yandex.com</td>
            </tr>
            <tr>
              <td>Type:</td>
              <td>server_cookie</td>
            </tr>
            <tr>
              <td>Expires in:</td>
              <td>11 months 30 days</td>
            </tr>
          </tbody>
        </table>
      </div>

      <strong>How can I control cookies on my browser?</strong>

      <p>
        As the means by which you can refuse cookies through your web browser
        controls vary from browser to browser, you should visit your browser's
        help menu for more information. The following is information about how
        to manage cookies on the most popular browsers:
      </p>

      <ul>
        <li>
          <a
            href="https://support.google.com/chrome/answer/95647#zippy=%2Callow-or-block-cookies"
            target="_blank"
            rel="noreferrer"
          >
            Chrome
          </a>
        </li>
        <li>
          <a
            href="https://support.microsoft.com/en-us/windows/manage-cookies-in-microsoft-edge-view-allow-block-delete-and-use-168dab11-0753-043d-7c16-ede5947fc64d"
            target="_blank"
            rel="noreferrer"
          >
            Internet Explorer
          </a>
        </li>
        <li>
          <a
            href="https://support.mozilla.org/en-US/kb/enhanced-tracking-protection-firefox-desktop?redirectslug=enable-and-disable-cookies-website-preferences&redirectlocale=en-US"
            target="_blank"
            rel="noreferrer"
          >
            Firefox
          </a>
        </li>
        <li>
          <a
            href="https://support.apple.com/en-ie/guide/safari/sfri11471/mac"
            target="_blank"
            rel="noreferrer"
          >
            Safari
          </a>
        </li>
        <li>
          <a
            href="https://support.microsoft.com/en-us/microsoft-edge/microsoft-edge-browsing-data-and-privacy-bb8174ba-9d73-dcf2-9b4a-c582b4e640dd"
            target="_blank"
            rel="noreferrer"
          >
            Edge
          </a>
        </li>
        <li>
          <a
            href="https://help.opera.com/en/latest/web-preferences/"
            target="_blank"
            rel="noreferrer"
          >
            Opera
          </a>
        </li>
      </ul>

      <p>
        In addition, most advertising networks offer you a way to opt out of
        targeted advertising. If you would like to find out more information,
        please visit:
      </p>

      <ul>
        <li>
          <a
            href="https://optout.aboutads.info/?c=2&lang=EN"
            target="_blank"
            rel="noreferrer"
          >
            Digital Advertising Alliance
          </a>
        </li>
        <li>
          <a href="https://youradchoices.ca/" target="_blank" rel="noreferrer">
            Digital Advertising Alliance of Canada
          </a>
        </li>
        <li>
          <a
            href="https://www.youronlinechoices.com/"
            target="_blank"
            rel="noreferrer"
          >
            European Interactive Digital Advertising Alliance
          </a>
        </li>
      </ul>

      <strong>What about other tracking technologies, like web beacons?</strong>

      <p>
        Cookies are not the only way to recognize or track visitors to a
        website. We may use other, similar technologies from time to time, like
        web beacons (sometimes called "tracking pixels" or "clear gifs"). These
        are tiny graphics files that contain a unique identifier that enables us
        to recognize when someone has visited our Website or opened an email
        including them. This allows us, for example, to monitor the traffic
        patterns of users from one page within a website to another, to deliver
        or communicate with cookies, to understand whether you have come to the
        website from an online advertisement displayed on a third-party website,
        to improve site performance, and to measure the success of email
        marketing campaigns. In many instances, these technologies are reliant
        on cookies to function properly, and so declining cookies will impair
        their functioning.
      </p>

      <strong>Do you use Flash cookies or Local Shared Objects?</strong>

      <p>
        Websites may also use so-called "Flash Cookies" (also known as Local
        Shared Objects or "LSOs") to, among other things, collect and store
        information about your use of our services, fraud prevention, and for
        other site operations.
      </p>

      <p>
        If you do not want Flash Cookies stored on your computer, you can adjust
        the settings of your Flash player to block Flash Cookies storage using
        the tools contained in the{' '}
        <a
          href="https://www.macromedia.com/support/documentation/en/flashplayer/help/settings_manager07.html"
          target="_blank"
          rel="noreferrer"
        >
          Website Storage Settings Panel
        </a>
        . You can also control Flash Cookies by going to the{' '}
        <a
          href="https://www.macromedia.com/support/documentation/en/flashplayer/help/settings_manager03.html"
          target="_blank"
          rel="noreferrer"
        >
          Global Storage Settings Panel
        </a>{' '}
        and following the instructions (which may include instructions that
        explain, for example, how to delete existing Flash Cookies (referred to
        "information" on the Macromedia site), how to prevent Flash LSOs from
        being placed on your computer without your being asked, and (for Flash
        Player 8 and later) how to block Flash Cookies that are not being
        delivered by the operator of the page you are on at the time).
      </p>

      <p>
        Please note that setting the Flash Player to restrict or limit
        acceptance of Flash Cookies may reduce or impede the functionality of
        some Flash applications, including, potentially, Flash applications used
        in connection with our services or online content.
      </p>

      <strong>Do you serve targeted advertising?</strong>

      <p>
        Third parties may serve cookies on your computer or mobile device to
        serve advertising through our Website. These companies may use
        information about your visits to this and other websites in order to
        provide relevant advertisements about goods and services that you may be
        interested in. They may also employ technology that is used to measure
        the effectiveness of advertisements. They can accomplish this by using
        cookies or web beacons to collect information about your visits to this
        and other sites in order to provide relevant advertisements about goods
        and services of potential interest to you. The information collected
        through this process does not enable us or them to identify your name,
        contact details, or other details that directly identify you unless you
        choose to provide these.
      </p>

      <strong>How often will you update this Cookie Policy?</strong>

      <p>
        We may update this Cookie Policy from time to time in order to reflect,
        for example, changes to the cookies we use or for other operational,
        legal, or regulatory reasons. Please therefore revisit this Cookie
        Policy regularly to stay informed about our use of cookies and related
        technologies.
      </p>

      <p>
        The date at the top of this Cookie Policy indicates when it was last
        updated.
      </p>

      <strong>Where can I get further information?</strong>

      <p>
        If you have any questions about our use of cookies or other
        technologies, please email us at{' '}
        <a href="mailto:support@hexlet.io">support@hexlet.io</a> or by post to:
        <ul className="fw-medium list-unstyled mt-3">
          <li>TOO "Hexlet"</li>
          <li>Auezov St 14 A</li>
          <li>Almaty, Almaty 050000</li>
          <li>Kazakhstan</li>
          <li>Phone: (+1)8 800 100 22 47</li>
          <li>Fax: (+1)8 800 100 22 47</li>
        </ul>
      </p>
    </>
  );
}
