import * as Routes from "@/routes.js";

export default function PrivacyEn() {
  return (
    <>
      <p className="fw-medium">Last updated March 28, 2025</p>

      <p>
        This Privacy Notice for TOO "Hexlet" ("we," "us," or "our"), describes
        how and why we might access, collect, store, use, and/or share
        ("process") your personal information when you use our services
        ("Services"), including when you:
        <ul>
          <li>
            Visit our website at{" "}
            <a href={Routes.root_path()} target="_blank" rel="noreferrer">
              https://code-basics.com
            </a>
            , or any website of ours that links to this Privacy Notice
          </li>
          <li>
            Engage with us in other related ways, including any sales,
            marketing, or events
          </li>
        </ul>
      </p>

      <p>
        <strong>Questions or concerns?</strong> Reading this Privacy Notice will
        help you understand your privacy rights and choices. We are responsible
        for making decisions about how your personal information is processed.
        If you do not agree with our policies and practices, please do not use
        our Services. If you still have any questions or concerns, please
        contact us at <a href="mailto:support@hexlet.io">support@hexlet.io</a>.
      </p>

      <h3>SUMMARY OF KEY POINTS</h3>

      <p className="fst-italic fw-medium">
        This summary provides key points from our Privacy Notice, but you can
        find out more details about any of these topics by clicking the link
        following each key point or by using our{" "}
        <a href="#table-of-contents">table of contents</a> below to find the
        section you are looking for.
      </p>

      <p>
        <strong>What personal information do we process? </strong> When you
        visit, use, or navigate our Services, we may process personal
        information depending on how you interact with us and the Services, the
        choices you make, and the products and features you use. Learn more
        about{" "}
        <a href="#personal-information">
          personal information you disclose to us
        </a>
        .
      </p>

      <p>
        <strong>Do we process any sensitive personal information? </strong> Some
        of the information may be considered "special" or "sensitive" in certain
        jurisdictions, for example your racial or ethnic origins, sexual
        orientation, and religious beliefs. We do not process sensitive personal
        information.
      </p>

      <p>
        <strong>Do we collect any information from third parties? </strong> We
        may collect information from public databases, marketing partners,
        social media platforms, and other outside sources. Learn more about{" "}
        <a href="#information-from-other-source">
          information collected from other sources
        </a>
        .
      </p>

      <p>
        <strong>
          In what situations and with which parties do we share personal
          information?
        </strong>{" "}
        We process your information to provide, improve, and administer our
        Services, communicate with you, for security and fraud prevention, and
        to comply with law. We may also process your information for other
        purposes with your consent. We process your information only when we
        have a valid legal reason to do so. Learn more about{" "}
        <a href="#process-information">how we process your information</a>.
      </p>

      <p>
        <strong>
          In what situations and with which parties do we share personal
          information?
        </strong>{" "}
        We may share information in specific situations and with specific third
        parties. Learn more about{" "}
        <a href="#share-information">
          when and with whom we share your personal information
        </a>
        .
      </p>

      <p>
        <strong>How do we keep your information safe?</strong> We have adequate
        organizational and technical processes and procedures in place to
        protect your personal information. However, no electronic transmission
        over the internet or information storage technology can be guaranteed to
        be 100% secure, so we cannot promise or guarantee that hackers,
        cybercriminals, or other unauthorized third parties will not be able to
        defeat our security and improperly collect, access, steal, or modify
        your information. Learn more about{" "}
        <a href="#save-information">how we keep your information safe</a>.
      </p>

      <p>
        <strong>What are your rights?</strong> Depending on where you are
        located geographically, the applicable privacy law may mean you have
        certain rights regarding your personal information. Learn more about{" "}
        <a href="#privacy-rights">your privacy rights</a>.
      </p>

      <p>
        <strong>How do you exercise your rights?</strong> The easiest way to
        exercise your rights is by submitting a{" "}
        <a href="https://app.termly.io/notify/7a5b7430-459b-4ae6-940d-56a91aed21ec">
          data subject access request
        </a>
        , or by contacting us. We will consider and act upon any request in
        accordance with applicable data protection laws.
      </p>

      <p>
        Want to learn more about what we do with any information we collect?{" "}
        <a href="#personal-information">Review the Privacy Notice in full.</a>.
      </p>

      <h3 id="table-of-contents">TABLE OF CONTENTS</h3>

      <ul className="list-unstyled">
        <li>
          <a href="#personal-information">1. WHAT INFORMATION DO WE COLLECT?</a>
        </li>
        <li>
          <a href="#process-information">
            2. HOW DO WE PROCESS YOUR INFORMATION?
          </a>
        </li>
        <li>
          <a href="#legal-bases">
            3. WHAT LEGAL BASES DO WE RELY ON TO PROCESS YOUR PERSONAL
            INFORMATION?
          </a>
        </li>
        <li>
          <a href="#share-information">
            4. WHEN AND WITH WHOM DO WE SHARE YOUR PERSONAL INFORMATION?
          </a>
        </li>
        <li>
          <a href="#use-cookies">
            5. DO WE USE COOKIES AND OTHER TRACKING TECHNOLOGIES?
          </a>
        </li>
        <li>
          <a href="#intelligence-based-products">
            6. DO WE OFFER ARTIFICIAL INTELLIGENCE-BASED PRODUCTS?
          </a>
        </li>
        <li>
          <a href="#data-storage">7. HOW LONG DO WE KEEP YOUR INFORMATION?</a>
        </li>
        <li>
          <a href="#save-information">
            8. HOW DO WE KEEP YOUR INFORMATION SAFE?
          </a>
        </li>
        <li>
          <a href="#minors-information">
            9. DO WE COLLECT INFORMATION FROM MINORS?
          </a>
        </li>
        <li>
          <a href="#privacy-rights">10. WHAT ARE YOUR PRIVACY RIGHTS?</a>
        </li>
        <li>
          <a href="#do-not-track-features">
            11. CONTROLS FOR DO-NOT-TRACK FEATURES
          </a>
        </li>
        <li>
          <a href="#specific-privacy-rights">
            12. DO UNITED STATES RESIDENTS HAVE SPECIFIC PRIVACY RIGHTS?
          </a>
        </li>
        <li>
          <a href="#specific-rights-of-other-regions">
            13. DO OTHER REGIONS HAVE SPECIFIC PRIVACY RIGHTS?
          </a>
        </li>
        <li>
          <a href="#notice-update">14. DO WE MAKE UPDATES TO THIS NOTICE?</a>
        </li>
        <li>
          <a href="#contacts">15. HOW CAN YOU CONTACT US ABOUT THIS NOTICE?</a>
        </li>
        <li>
          <a href="#change-information">
            16. HOW CAN YOU REVIEW, UPDATE, OR DELETE THE DATA WE COLLECT FROM
            YOU?
          </a>
        </li>
      </ul>

      <h3 id="personal-information">1. WHAT INFORMATION DO WE COLLECT?</h3>

      <h4>Personal information you disclose to us</h4>

      <p className="fst-italic">
        <strong>In Short:</strong> We collect personal information that you
        provide to us.
      </p>

      <p>
        We collect personal information that you voluntarily provide to us when
        you register on the Services, express an interest in obtaining
        information about us or our products and Services, when you participate
        in activities on the Services, or otherwise when you contact us.
      </p>

      <p>
        <strong>Personal Information Provided by You. </strong> The personal
        information that we collect depends on the context of your interactions
        with us and the Services, the choices you make, and the products and
        features you use. The personal information we collect may include the
        following:
        <ul>
          <li>names</li>
          <li>email addresses</li>
          <li>passwords</li>
        </ul>
      </p>

      <p>
        <strong>Sensitive Information.</strong> We do not process sensitive
        information.
      </p>

      <p>
        All personal information that you provide to us must be true, complete,
        and accurate, and you must notify us of any changes to such personal
        information.
      </p>

      <h4>Information automatically collected</h4>

      <p className="fst-italic">
        <strong>In Short:</strong> Some information — such as your Internet
        Protocol (IP) address and/or browser and device characteristics — is
        collected automatically when you visit our Services.
      </p>

      <p>
        We automatically collect certain information when you visit, use, or
        navigate the Services. This information does not reveal your specific
        identity (like your name or contact information) but may include device
        and usage information, such as your IP address, browser and device
        characteristics, operating system, language preferences, referring URLs,
        device name, country, location, information about how and when you use
        our Services, and other technical information. This information is
        primarily needed to maintain the security and operation of our Services,
        and for our internal analytics and reporting purposes.
      </p>

      <p>
        Like many businesses, we also collect information through cookies and
        similar technologies. You can find out more about this in our Cookie
        Notice:{" "}
        <a
          href={Routes.page_path("cookie_policy")}
          target="_blank"
          rel="noreferrer"
        >
          https://code-basics.com/pages/cookie_policy
        </a>
        .
      </p>

      <p>
        The information we collect includes:
        <ul>
          <li>
            Log and Usage Data. Log and usage data is service-related,
            diagnostic, usage, and performance information our servers
            automatically collect when you access or use our Services and which
            we record in log files. Depending on how you interact with us, this
            log data may include your IP address, device information, browser
            type, and settings and information about your activity in the
            Services (such as the date/time stamps associated with your usage,
            pages and files viewed, searches, and other actions you take such as
            which features you use), device event information (such as system
            activity, error reports (sometimes called "crash dumps"), and
            hardware settings).
          </li>
          <li>
            Device Data. We collect device data such as information about your
            computer, phone, tablet, or other device you use to access the
            Services. Depending on the device used, this device data may include
            information such as your IP address (or proxy server), device and
            application identification numbers, location, browser type, hardware
            model, Internet service provider and/or mobile carrier, operating
            system, and system configuration information.
          </li>
          <li>
            Location Data. We collect location data such as information about
            your device's location, which can be either precise or imprecise.
            How much information we collect depends on the type and settings of
            the device you use to access the Services. For example, we may use
            GPS and other technologies to collect geolocation data that tells us
            your current location (based on your IP address). You can opt out of
            allowing us to collect this information either by refusing access to
            the information or by disabling your Location setting on your
            device. However, if you choose to opt out, you may not be able to
            use certain aspects of the Services.
          </li>
        </ul>
      </p>

      <h4>Information collected from other sources</h4>

      <p className="fst-italic">
        <strong>In Short:</strong> We may collect limited data from public
        databases, marketing partners, and other outside sources.
      </p>

      <p>
        In order to enhance our ability to provide relevant marketing, offers,
        and services to you and update our records, we may obtain information
        about you from other sources, such as public databases, joint marketing
        partners, affiliate programs, data providers, and from other third
        parties. This information includes mailing addresses, job titles, email
        addresses, phone numbers, intent data (or user behavior data), Internet
        Protocol (IP) addresses, social media profiles, social media URLs, and
        custom profiles, for purposes of targeted advertising and event
        promotion.
      </p>

      <h3 id="process-information">2. HOW DO WE PROCESS YOUR INFORMATION?</h3>

      <p className="fst-italic">
        <strong>In Short:</strong> We process your information to provide,
        improve, and administer our Services, communicate with you, for security
        and fraud prevention, and to comply with law. We may also process your
        information for other purposes with your consent.
      </p>

      <p className="fw-medium">
        We process your personal information for a variety of reasons, depending
        on how you interact with our Services, including:
      </p>
      <ul>
        <li>
          <span className="fw-medium">
            To facilitate account creation and authentication and otherwise
            manage user accounts.
          </span>{" "}
          We may process your information so you can create and log in to your
          account, as well as keep your account in working order.
        </li>
        <li>
          <span className="fw-medium">
            To save or protect an individual's vital interest.
          </span>{" "}
          We may process your information when necessary to save or protect an
          individual’s vital interest, such as to prevent harm.
        </li>
      </ul>

      <h3 id="legal-bases">
        3. WHAT LEGAL BASES DO WE RELY ON TO PROCESS YOUR INFORMATION?
      </h3>

      <p className="fst-italic">
        <strong>In Short: </strong>We only process your personal information
        when we believe it is necessary and we have a valid legal reason (i.e.,
        legal basis) to do so under applicable law, like with your consent, to
        comply with laws, to provide you with services to enter into or fulfill
        our contractual obligations, to protect your rights, or to fulfill our
        legitimate business interests.
      </p>

      <p className="fw-medium text-decoration-underline">
        If you are located in the EU or UK, this section applies to you.
      </p>

      <p>
        The General Data Protection Regulation (GDPR) and UK GDPR require us to
        explain the valid legal bases we rely on in order to process your
        personal information. As such, we may rely on the following legal bases
        to process your personal information:
      </p>

      <ul>
        <li>
          <span className="fw-medium">Consent.</span> We may process your
          information if you have given us permission (i.e., consent) to use
          your personal information for a specific purpose. You can withdraw
          your consent at any time. Learn more about
          <a href="#withdrawing-consent"> withdrawing your consent.</a>
        </li>
        <li>
          <span className="fw-medium">Legal Obligations.</span> We may process
          your information where we believe it is necessary for compliance with
          our legal obligations, such as to cooperate with a law enforcement
          body or regulatory agency, exercise or defend our legal rights, or
          disclose your information as evidence in litigation in which we are
          involved.
        </li>
        <li>
          <span className="fw-medium">Vital Interests.</span> We may process
          your information where we believe it is necessary to protect your
          vital interests or the vital interests of a third party, such as
          situations involving potential threats to the safety of any person.
        </li>
      </ul>

      <p className="fw-medium text-decoration-underline">
        If you are located in Canada, this section applies to you.
      </p>

      <p>
        We may process your information if you have given us specific permission
        (i.e., express consent) to use your personal information for a specific
        purpose, or in situations where your permission can be inferred (i.e.,
        implied consent). You can
        <a href="#withdrawing-consent"> withdrawing your consent</a> at any
        time.
      </p>

      <p>
        In some exceptional cases, we may be legally permitted under applicable
        law to process your information without your consent, including, for
        example:
        <ul>
          <li>
            If collection is clearly in the interests of an individual and
            consent cannot be obtained in a timely way
          </li>
          <li>For investigations and fraud detection and prevention</li>
          <li>For business transactions provided certain conditions are met</li>
          <li>
            If it is contained in a witness statement and the collection is
            necessary to assess, process, or settle an insurance claim
          </li>
          <li>
            For identifying injured, ill, or deceased persons and communicating
            with next of kin
          </li>
          <li>
            If we have reasonable grounds to believe an individual has been, is,
            or may be victim of financial abuse
          </li>
          <li>
            If it is reasonable to expect collection and use with consent would
            compromise the availability or the accuracy of the information and
            the collection is reasonable for purposes related to investigating a
            breach of an agreement or a contravention of the laws of Canada or a
            province
          </li>
          <li>
            If disclosure is required to comply with a subpoena, warrant, court
            order, or rules of the court relating to the production of records
          </li>
          <li>
            If it was produced by an individual in the course of their
            employment, business, or profession and the collection is consistent
            with the purposes for which the information was produced
          </li>
          <li>
            If the collection is solely for journalistic, artistic, or literary
            purposes
          </li>
          <li>
            If the information is publicly available and is specified by the
            regulations
          </li>
        </ul>
      </p>

      <h3 id="share-information">
        4. WHEN AND WITH WHOM DO WE SHARE YOUR PERSONAL INFORMATION?
      </h3>

      <p className="fst-italic">
        <strong>In Short: </strong>We may share information in specific
        situations described in this section and/or with the following third
        parties.
      </p>

      <p>
        We may need to share your personal information in the following
        situations:
        <ul>
          <li>
            <strong>Business Transfers.</strong> We may share or transfer your
            information in connection with, or during negotiations of, any
            merger, sale of company assets, financing, or acquisition of all or
            a portion of our business to another company.
          </li>
          <li>
            <strong>Affiliates.</strong> We may share your information with our
            affiliates, in which case we will require those affiliates to honor
            this Privacy Notice. Affiliates include our parent company and any
            subsidiaries, joint venture partners, or other companies that we
            control or that are under common control with us.
          </li>
        </ul>
      </p>

      <h3 id="use-cookies">
        5. DO WE USE COOKIES AND OTHER TRACKING TECHNOLOGIES?
      </h3>

      <p className="fst-italic">
        <strong>In Short: </strong> We may use cookies and other tracking
        technologies to collect and store your information.
      </p>

      <p>
        We may use cookies and similar tracking technologies (like web beacons
        and pixels) to gather information when you interact with our Services.
        Some online tracking technologies help us maintain the security of our
        Services and your account, prevent crashes, fix bugs, save your
        preferences, and assist with basic site functions
      </p>

      <p>
        We also permit third parties and service providers to use online
        tracking technologies on our Services for analytics and advertising,
        including to help manage and display advertisements, to tailor
        advertisements to your interests, or to send abandoned shopping cart
        reminders (depending on your communication preferences). The third
        parties and service providers use their technology to provide
        advertising about products and services tailored to your interests which
        may appear either on our Services or on other websites.
      </p>

      <p>
        To the extent these online tracking technologies are deemed to be a
        "sale"/"sharing" (which includes targeted advertising, as defined under
        the applicable laws) under applicable US state laws, you can opt out of
        these online tracking technologies by submitting a request as described
        below under section{" "}
        <a href="#specific-privacy-rights">
          DO UNITED STATES RESIDENTS HAVE SPECIFIC PRIVACY RIGHTS?"
        </a>
      </p>

      <p>
        Specific information about how we use such technologies and how you can
        refuse certain cookies is set out in our Cookie Notice:{" "}
        <a
          href={Routes.page_path("cookie_policy")}
          target="_blank"
          rel="noreferrer"
        >
          https://code-basics.com/pages/cookie_policy
        </a>
      </p>

      <h4>Google Analytics</h4>

      <p>
        We may share your information with Google Analytics to track and analyze
        the use of the Services. The Google Analytics Advertising Features that
        we may use include: Remarketing with Google Analytics. To opt out of
        being tracked by Google Analytics across the Services, visit{" "}
        <a
          href="https://tools.google.com/dlpage/gaoptout"
          target="_blank"
          rel="noreferrer"
        >
          https://tools.google.com/dlpage/gaoptout
        </a>
        . You can opt out of Google Analytics Advertising Features through{" "}
        <a
          href="https://adssettings.google.com/"
          target="_blank"
          rel="noreferrer"
        >
          Ads Settings
        </a>{" "}
        and Ad Settings for mobile apps. Other opt out means include{" "}
        <a
          href="https://optout.networkadvertising.org/?c=1"
          target="_blank"
          rel="noreferrer"
        >
          http://optout.networkadvertising.org/
        </a>{" "}
        and{" "}
        <a
          href="http://www.networkadvertising.org/mobile-choice"
          target="_blank"
          rel="noreferrer"
        >
          http://www.networkadvertising.org/mobile-choice
        </a>
        . For more information on the privacy practices of Google, please visit
        the{" "}
        <a
          href="https://policies.google.com/privacy"
          target="_blank"
          rel="noreferrer"
        >
          Google Privacy & Terms page.
        </a>
      </p>

      <h3 id="intelligence-based-products">
        6. DO WE OFFER ARTIFICIAL INTELLIGENCE-BASED PRODUCTS?
      </h3>

      <p className="fst-italic">
        <strong>In Short: </strong>We offer products, features, or tools powered
        by artificial intelligence, machine learning, or similar technologies.
      </p>

      <p>
        As part of our Services, we offer products, features, or tools powered
        by artificial intelligence, machine learning, or similar technologies
        (collectively, "AI Products"). These tools are designed to enhance your
        experience and provide you with innovative solutions. The terms in this
        Privacy Notice govern your use of the AI Products within our Services.
      </p>

      <strong>Use of AI Technologies</strong>

      <p>
        We provide the AI Products through third-party service providers ("AI
        Service Providers"), including OpenAI. As outlined in this Privacy
        Notice, your input, output, and personal information will be shared with
        and processed by these AI Service Providers to enable your use of our AI
        Products for purposes outlined in{" "}
        <a href="#legal-bases">
          "WHAT LEGAL BASES DO WE RELY ON TO PROCESS YOUR PERSONAL INFORMATION?"
        </a>
        You must not use the AI Products in any way that violates the terms or
        policies of any AI Service Provider.
      </p>

      <strong>Our AI Products</strong>

      <p>
        Our AI Products are designed for the following functions:
        <ul>
          <li>AI applications</li>
        </ul>
      </p>

      <strong>How We Process Your Data Using AI</strong>

      <p>
        All personal information processed using our AI Products is handled in
        line with our Privacy Notice and our agreement with third parties. This
        ensures high security and safeguards your personal information
        throughout the process, giving you peace of mind about your data's
        safety.
      </p>

      <h3 id="save-information">7. HOW LONG DO WE KEEP YOUR INFORMATION?</h3>

      <p className="fst-italic">
        <strong>In Short:</strong> We keep your information for as long as
        necessary to fulfill the purposes outlined in this Privacy Notice unless
        otherwise required by law.
      </p>

      <p>
        We will only keep your personal information for as long as it is
        necessary for the purposes set out in this Privacy Notice, unless a
        longer retention period is required or permitted by law (such as tax,
        accounting, or other legal requirements). No purpose in this notice will
        require us keeping your personal information for longer than the period
        of time in which users have an account with us.
      </p>

      <p>
        When we have no ongoing legitimate business need to process your
        personal information, we will either delete or anonymize such
        information, or, if this is not possible (for example, because your
        personal information has been stored in backup archives), then we will
        securely store your personal information and isolate it from any further
        processing until deletion is possible.
      </p>

      <h3 id="save-information">8. HOW DO WE KEEP YOUR INFORMATION SAFE?</h3>

      <p className="fst-italic">
        <strong>In Short: </strong>We aim to protect your personal information
        through a system of organizational and technical security measures.
      </p>

      <p>
        We have implemented appropriate and reasonable technical and
        organizational security measures designed to protect the security of any
        personal information we process. However, despite our safeguards and
        efforts to secure your information, no electronic transmission over the
        Internet or information storage technology can be guaranteed to be 100%
        secure, so we cannot promise or guarantee that hackers, cybercriminals,
        or other unauthorized third parties will not be able to defeat our
        security and improperly collect, access, steal, or modify your
        information. Although we will do our best to protect your personal
        information, transmission of personal information to and from our
        Services is at your own risk. You should only access the Services within
        a secure environment.
      </p>

      <h3 id="minors-information">9. DO WE COLLECT INFORMATION FROM MINORS?</h3>

      <p className="fst-italic">
        <strong>In Short: </strong>We do not knowingly collect data from or
        market to children under 18 years of age.
      </p>

      <p>
        We do not knowingly collect, solicit data from, or market to children
        under 18 years of age, nor do we knowingly sell such personal
        information. By using the Services, you represent that you are at least
        18 or that you are the parent or guardian of such a minor and consent to
        such minor dependent’s use of the Services. If we learn that personal
        information from users less than 18 years of age has been collected, we
        will deactivate the account and take reasonable measures to promptly
        delete such data from our records. If you become aware of any data we
        may have collected from children under age 18, please contact us at
        support@hexlet.io.
      </p>

      <h3 id="privacy-rights">10. WHAT ARE YOUR PRIVACY RIGHTS?</h3>

      <p className="fst-italic">
        <strong>In Short: </strong>Depending on your state of residence in the
        US or in some regions, such as the European Economic Area (EEA), United
        Kingdom (UK), Switzerland, and Canada, you have rights that allow you
        greater access to and control over your personal information. You may
        review, change, or terminate your account at any time, depending on your
        country, province, or state of residence.
      </p>

      <p>
        In some regions (like the EEA, UK, Switzerland, and Canada), you have
        certain rights under applicable data protection laws. These may include
        the right (i) to request access and obtain a copy of your personal
        information, (ii) to request rectification or erasure; (iii) to restrict
        the processing of your personal information; (iv) if applicable, to data
        portability; and (v) not to be subject to automated decision-making. In
        certain circumstances, you may also have the right to object to the
        processing of your personal information. You can make such a request by
        contacting us by using the contact details provided in the section{" "}
        <a href="#contacts">"HOW CAN YOU CONTACT US ABOUT THIS NOTICE?"</a>{" "}
        below.
      </p>

      <p>
        We will consider and act upon any request in accordance with applicable
        data protection laws.
      </p>

      <p>
        If you are located in the EEA or UK and you believe we are unlawfully
        processing your personal information, you also have the right to
        complain to your{" "}
        <a
          href="https://ec.europa.eu/newsroom/article29/items/612080"
          target="_blank"
          rel="noreferrer"
        >
          Member State data protection authority
        </a>{" "}
        or
        <a
          href="https://ico.org.uk/make-a-complaint/data-protection-complaints/"
          target="_blank"
          rel="noreferrer"
        >
          UK data protection authority
        </a>
        .
      </p>

      <p>
        If you are located in Switzerland, you may contact the{" "}
        <a
          href="https://www.edoeb.admin.ch/en"
          target="_blank"
          rel="noreferrer"
        >
          Federal Data Protection and Information Commissioner
        </a>
        .
      </p>

      <p>
        <strong id="withdrawing-consent">Withdrawing your consent: </strong>If
        we are relying on your consent to process your personal information,
        which may be express and/or implied consent depending on the applicable
        law, you have the right to withdraw your consent at any time. You can
        withdraw your consent at any time by contacting us by using the contact
        details provided in the section{" "}
        <a href="#contacts">"HOW CAN YOU CONTACT US ABOUT THIS NOTICE?"</a>{" "}
        below.
      </p>

      <p>
        However, please note that this will not affect the lawfulness of the
        processing before its withdrawal nor, when applicable law allows, will
        it affect the processing of your personal information conducted in
        reliance on lawful processing grounds other than consent.
      </p>

      <p>
        <strong>Opting out of marketing and promotional communications:</strong>{" "}
        You can unsubscribe from our marketing and promotional communications at
        any time by clicking on the unsubscribe link in the emails that we send,
        or by contacting us using the details provided in the section{" "}
        <a href="#contacts">"HOW CAN YOU CONTACT US ABOUT THIS NOTICE?"</a>{" "}
        below. You will then be removed from the marketing lists. However, we
        may still communicate with you — for example, to send you
        service-related messages that are necessary for the administration and
        use of your account, to respond to service requests, or for other
        non-marketing purposes.
      </p>

      <h4>Account Information</h4>

      <p>
        If you would at any time like to review or change the information in
        your account or terminate your account, you can:
        <ul>
          <li>Log in to your account settings and update your user account.</li>
        </ul>
      </p>

      <p>
        Upon your request to terminate your account, we will deactivate or
        delete your account and information from our active databases. However,
        we may retain some information in our files to prevent fraud,
        troubleshoot problems, assist with any investigations, enforce our legal
        terms and/or comply with applicable legal requirements.
      </p>

      <p>
        <strong>Cookies and similar technologies:</strong> Most Web browsers are
        set to accept cookies by default. If you prefer, you can usually choose
        to set your browser to remove cookies and to reject cookies. If you
        choose to remove cookies or reject cookies, this could affect certain
        features or services of our Services. For further information, please
        see our Cookie Notice:{" "}
        <a
          href={Routes.page_path("cookie_policy")}
          target="_blank"
          rel="noreferrer"
        >
          https://code-basics.com/pages/cookie_policy
        </a>
        .
      </p>

      <p>
        If you have questions or comments about your privacy rights, you may
        email us at support@hexlet.io.
      </p>

      <h3 id="do-not-track-features">11. CONTROLS FOR DO-NOT-TRACK FEATURES</h3>

      <p>
        Most web browsers and some mobile operating systems and mobile
        applications include a Do-Not-Track ("DNT") feature or setting you can
        activate to signal your privacy preference not to have data about your
        online browsing activities monitored and collected. At this stage, no
        uniform technology standard for recognizing and implementing DNT signals
        has been finalized. As such, we do not currently respond to DNT browser
        signals or any other mechanism that automatically communicates your
        choice not to be tracked online. If a standard for online tracking is
        adopted that we must follow in the future, we will inform you about that
        practice in a revised version of this Privacy Notice.
      </p>

      <p>
        California law requires us to let you know how we respond to web browser
        DNT signals. Because there currently is not an industry or legal
        standard for recognizing or honoring DNT signals, we do not respond to
        them at this time.
      </p>

      <h3 id="specific-privacy-rights">
        12. DO UNITED STATES RESIDENTS HAVE SPECIFIC PRIVACY RIGHTS?
      </h3>

      <p className="fst-italic">
        <strong>In Short: </strong>If you are a resident of California,
        Colorado, Connecticut, Delaware, Florida, Indiana, Iowa, Kentucky,
        Maryland, Minnesota, Montana, Nebraska, New Hampshire, New Jersey,
        Oregon, Rhode Island, Tennessee, Texas, Utah, or Virginia, you may have
        the right to request access to and receive details about the personal
        information we maintain about you and how we have processed it, correct
        inaccuracies, get a copy of, or delete your personal information. You
        may also have the right to withdraw your consent to our processing of
        your personal information. These rights may be limited in some
        circumstances by applicable law. More information is provided below.
      </p>

      <h4>Categories of Personal Information We Collect</h4>

      <p>
        We have collected the following categories of personal information in
        the past twelve (12) months:
      </p>

      <div className="table-responsive">
        <table className="table table-bordered">
          <tbody>
            <tr>
              <th className="text-center">Category</th>
              <th className="text-center">Examples</th>
              <th className="text-center">Collected</th>
            </tr>
            <tr>
              <td>A. Identifiers</td>
              <td>
                Contact details, such as real name, alias, postal address,
                telephone or mobile contact number, unique personal identifier,
                online identifier, Internet Protocol address, email address, and
                account name
              </td>
              <td>YES</td>
            </tr>
            <tr>
              <td>
                B. Personal information as defined in the California Customer
                Records statute
              </td>
              <td>
                Name, contact information, education, employment, employment
                history, and financial information
              </td>
              <td>YES</td>
            </tr>
            <tr>
              <td>
                C. Protected classification characteristics under state or
                federal law
              </td>
              <td>
                Gender, age, date of birth, race and ethnicity, national origin,
                marital status, and other demographic data
              </td>
              <td>NO</td>
            </tr>
            <tr>
              <td>D. Commercial information</td>
              <td>
                Transaction information, purchase history, financial details,
                and payment information
              </td>
              <td>NO</td>
            </tr>
            <tr>
              <td>E. Biometric information</td>
              <td>Fingerprints and voiceprints </td>
              <td>NO</td>
            </tr>
            <tr>
              <td>F. Internet or other similar network activity</td>
              <td>
                Browsing history, search history, online behavior, interest
                data, and interactions with our and other websites,
                applications, systems, and advertisements
              </td>
              <td>YES</td>
            </tr>
            <tr>
              <td>G. Geolocation data </td>
              <td>Device location </td>
              <td>YES</td>
            </tr>
            <tr>
              <td>H. Audio, electronic, sensory, or similar information</td>
              <td>
                Images and audio, video or call recordings created in connection
                with our business activities
              </td>
              <td>NO</td>
            </tr>
            <tr>
              <td>I. Professional or employment-related information</td>
              <td>
                Business contact details in order to provide you our Services at
                a business level or job title, work history, and professional
                qualifications if you apply for a job with us
              </td>
              <td>NO</td>
            </tr>
            <tr>
              <td>J. Education Information</td>
              <td>Student records and directory information </td>
              <td>NO</td>
            </tr>
            <tr>
              <td>K. Inferences drawn from collected personal information</td>
              <td>
                Inferences drawn from any of the collected personal information
                listed above to create a profile or summary about, for example,
                an individual’s preferences and characteristics
              </td>
              <td>YES</td>
            </tr>
            <tr>
              <td>L. Sensitive personal Information </td>
              <td />
              <td>NO</td>
            </tr>
          </tbody>
        </table>
      </div>

      <p>
        We may also collect other personal information outside of these
        categories through instances where you interact with us in person,
        online, or by phone or mail in the context of:
        <ul>
          <li>Receiving help through our customer support channels;</li>
          <li>Participation in customer surveys or contests; and</li>
          <li>
            Facilitation in the delivery of our Services and to respond to your
            inquiries.
          </li>
        </ul>
      </p>

      <p>
        We will use and retain the collected personal information as needed to
        provide the Services or for:
        <ul>
          <li>Category A - As long as the user has an account with us</li>
          <li>Category B - As long as the user has an account with us</li>
          <li>Category F - As long as the user has an account with us</li>
          <li>Category G - As long as the user has an account with us</li>
          <li>Category K - As long as the user has an account with us</li>
        </ul>
      </p>

      <h4>Sources of Personal Information</h4>

      <p>
        Learn more about the sources of personal information we collect in{" "}
        <a href="#personal-information">"WHAT INFORMATION DO WE COLLECT?"</a>
      </p>

      <h4>How We Use and Share Personal Information</h4>

      <p>
        Learn more about how we use your personal information in the section,{" "}
        <a href="#process-information">"HOW DO WE PROCESS YOUR INFORMATION?"</a>
      </p>

      <strong>Will your information be shared with anyone else?</strong>

      <p>
        We may disclose your personal information with our service providers
        pursuant to a written contract between us and each service provider.
        Learn more about how we disclose personal information to in the section,{" "}
        <a href="#share-information">
          "WHEN AND WITH WHOM DO WE SHARE YOUR PERSONAL INFORMATION?"
        </a>
      </p>

      <p>
        We may use your personal information for our own business purposes, such
        as for undertaking internal research for technological development and
        demonstration. This is not considered to be "selling" of your personal
        information.
      </p>

      <p>
        We have not disclosed, sold, or shared any personal information to third
        parties for a business or commercial purpose in the preceding twelve
        (12) months. We will not sell or share personal information in the
        future belonging to website visitors, users, and other consumers.
      </p>

      <h4>Your Rights</h4>

      <p>
        You have rights under certain US state data protection laws. However,
        these rights are not absolute, and in certain cases, we may decline your
        request as permitted by law. These rights include:
        <ul>
          <li>
            Right to know whether or not we are processing your personal data
          </li>
          <li>Right to access your personal data</li>
          <li>Right to correct inaccuracies in your personal data</li>
          <li>Right to request the deletion of your personal data</li>
          <li>
            Right to obtain a copy of the personal data you previously shared
            with us
          </li>
          <li>Right to non-discrimination for exercising your rights</li>
          <li>
            Right to opt out of the processing of your personal data if it is
            used for targeted advertising (or sharing as defined under
            California’s privacy law), the sale of personal data, or profiling
            in furtherance of decisions that produce legal or similarly
            significant effects ("profiling")
          </li>
        </ul>
      </p>

      <p>
        Depending upon the state where you live, you may also have the following
        rights:
        <ul>
          <li>
            Right to access the categories of personal data being processed (as
            permitted by applicable law, including the privacy law in Minnesota)
          </li>
          <li>
            Right to obtain a list of the categories of third parties to which
            we have disclosed personal data (as permitted by applicable law,
            including the privacy law in California, Delaware, and Maryland)
          </li>
          <li>
            Right to obtain a list of specific third parties to which we have
            disclosed personal data (as permitted by applicable law, including
            the privacy law in Minnesota and Oregon)
          </li>
          <li>
            Right to review, understand, question, and correct how personal data
            has been profiled (as permitted by applicable law, including the
            privacy law in Minnesota)
          </li>
          <li>
            Right to limit use and disclosure of sensitive personal data (as
            permitted by applicable law, including the privacy law in
            California)
          </li>
          <li>
            Right to opt out of the collection of sensitive data and personal
            data collected through the operation of a voice or facial
            recognition feature (as permitted by applicable law, including the
            privacy law in Florida)
          </li>
        </ul>
      </p>

      <h4>How to Exercise Your Rights</h4>

      <p>
        To exercise these rights, you can contact us by submitting a{" "}
        <a
          href="https://app.termly.io/notify/7a5b7430-459b-4ae6-940d-56a91aed21ec"
          target="_blank"
          rel="noreferrer"
        >
          data subject access request
        </a>
        , by emailing us at support@hexlet.io, or by referring to the contact
        details at the bottom of this document.
      </p>

      <p>
        Under certain US state data protection laws, you can designate an
        authorized agent to make a request on your behalf. We may deny a request
        from an authorized agent that does not submit proof that they have been
        validly authorized to act on your behalf in accordance with applicable
        laws.
      </p>

      <h4>Request Verification</h4>

      <p>
        Upon receiving your request, we will need to verify your identity to
        determine you are the same person about whom we have the information in
        our system. We will only use personal information provided in your
        request to verify your identity or authority to make the request.
        However, if we cannot verify your identity from the information already
        maintained by us, we may request that you provide additional information
        for the purposes of verifying your identity and for security or
        fraud-prevention purposes.
      </p>

      <p>
        If you submit the request through an authorized agent, we may need to
        collect additional information to verify your identity before processing
        your request and the agent will need to provide a written and signed
        permission from you to submit such request on your behalf.
      </p>

      <h4>Appeals</h4>

      <p>
        Under certain US state data protection laws, if we decline to take
        action regarding your request, you may appeal our decision by emailing
        us at support@hexlet.io. We will inform you in writing of any action
        taken or not taken in response to the appeal, including a written
        explanation of the reasons for the decisions. If your appeal is denied,
        you may submit a complaint to your state attorney general.
      </p>

      <h4>California "Shine The Light" Law</h4>

      <p>
        California Civil Code Section 1798.83, also known as the "Shine The
        Light" law, permits our users who are California residents to request
        and obtain from us, once a year and free of charge, information about
        categories of personal information (if any) we disclosed to third
        parties for direct marketing purposes and the names and addresses of all
        third parties with which we shared personal information in the
        immediately preceding calendar year. If you are a California resident
        and would like to make such a request, please submit your request in
        writing to us by using the contact details provided in the section{" "}
        <a href="#contacts">"HOW CAN YOU CONTACT US ABOUT THIS NOTICE?"</a>.
      </p>

      <h3 id="specific-rights-of-other-regions">
        13. DO OTHER REGIONS HAVE SPECIFIC PRIVACY RIGHTS?
      </h3>

      <p className="fst-italic">
        <strong>In Short: </strong>You may have additional rights based on the
        country you reside in.
      </p>

      <h4>Australia and New Zealand</h4>

      <p>
        We collect and process your personal information under the obligations
        and conditions set by Australia's Privacy Act 1988 and New Zealand's
        Privacy Act 2020 (Privacy Act).
      </p>

      <p>
        This Privacy Notice satisfies the notice requirements defined in both
        Privacy Acts, in particular: what personal information we collect from
        you, from which sources, for which purposes, and other recipients of
        your personal information.
      </p>

      <p>
        If you do not wish to provide the personal information necessary to
        fulfill their applicable purpose, it may affect our ability to provide
        our services, in particular:
        <ul>
          <li>offer you the products or services that you want</li>
          <li>respond to or help with your requests</li>
          <li>manage your account with us</li>
          <li>confirm your identity and protect your account</li>
        </ul>
      </p>

      <p>
        At any time, you have the right to request access to or correction of
        your personal information. You can make such a request by contacting us
        by using the contact details provided in the section{" "}
        <a href="#change-information">
          "HOW CAN YOU REVIEW, UPDATE, OR DELETE THE DATA WE COLLECT FROM YOU?"
        </a>
        .
      </p>

      <p>
        If you believe we are unlawfully processing your personal information,
        you have the right to submit a complaint about a breach of the
        Australian Privacy Principles to the{" "}
        <a
          href="https://www.oaic.gov.au/privacy/privacy-complaints/lodge-a-privacy-complaint-with-us"
          target="_blank"
          rel="noreferrer"
        >
          Office of the Australian Information Commissioner
        </a>
        and a breach of New Zealand's Privacy Principles to the{" "}
        <a
          href="https://www.privacy.org.nz/your-rights/making-a-complaint/"
          target="_blank"
          rel="noreferrer"
        >
          Office of New Zealand Privacy Commissioner
        </a>
        .
      </p>

      <h4>Republic of South Africa</h4>

      <p>
        At any time, you have the right to request access to or correction of
        your personal information. You can make such a request by contacting us
        by using the contact details provided in the section{" "}
        <a href="#change-information">
          "HOW CAN YOU REVIEW, UPDATE, OR DELETE THE DATA WE COLLECT FROM YOU?"
        </a>
        .
      </p>

      <p>
        If you are unsatisfied with the manner in which we address any complaint
        with regard to our processing of personal information, you can contact
        the office of the regulator, the details of which are:
        <ul className="list-unstyled">
          <li>
            <a
              href="https://inforegulator.org.za/"
              target="_blank"
              rel="noreferrer"
            >
              The Information Regulator (South Africa)
            </a>
          </li>
          <li>
            General enquiries:
            <a href="mailto:enquiries@inforegulator.org.za">
              enquiries@inforegulator.org.za
            </a>
          </li>
          <li>
            Complaints complete POPIA/PAIA form
            <a href="mailto:PAIAComplaints@inforegulator.org.za">
              PAIAComplaints@inforegulator.org.za
            </a>{" "}
            &{" "}
            <a href="mailto:POPIAComplaints@inforegulator.org.za">
              POPIAComplaints@inforegulator.org.za
            </a>
          </li>
        </ul>
      </p>

      <h3 id="notice-update">14. DO WE MAKE UPDATES TO THIS NOTICE?</h3>

      <p className="fst-italic">
        <strong>In Short: </strong>Yes, we will update this notice as necessary
        to stay compliant with relevant laws
      </p>

      <p>
        We may update this Privacy Notice from time to time. The updated version
        will be indicated by an updated "Revised" date at the top of this
        Privacy Notice. If we make material changes to this Privacy Notice, we
        may notify you either by prominently posting a notice of such changes or
        by directly sending you a notification. We encourage you to review this
        Privacy Notice frequently to be informed of how we are protecting your
        information.
      </p>

      <h3 id="contacts">15. HOW CAN YOU CONTACT US ABOUT THIS NOTICE?</h3>

      <p>
        If you have questions or comments about this notice, you may email us at{" "}
        <a href="mailto:support@hexlet.io">support@hexlet.io</a> or contact us
        by post at:
        <ul className="list-unstyled">
          <li>TOO "Hexlet"</li>
          <li>Auezov St 14 A</li>
          <li>Almaty, Almaty 050000</li>
          <li>Kazakhstan</li>
        </ul>
      </p>

      <h3 id="change-information">
        16. HOW CAN YOU REVIEW, UPDATE, OR DELETE THE DATA WE COLLECT FROM YOU?
      </h3>

      <p>
        Based on the applicable laws of your country or state of residence in
        the US, you may have the right to request access to the personal
        information we collect from you, details about how we have processed it,
        correct inaccuracies, or delete your personal information. You may also
        have the right to withdraw your consent to our processing of your
        personal information. These rights may be limited in some circumstances
        by applicable law. To request to review, update, or delete your personal
        information, please fill out and submit a{" "}
        <a
          href="https://app.termly.io/notify/7a5b7430-459b-4ae6-940d-56a91aed21ec"
          target="_blank"
          rel="noreferrer"
        >
          {" "}
          data subject access request
        </a>
        .
      </p>

      <ul className="fw-medium list-unstyled">
        <li>TOO "Hexlet"</li>
        <li>Auezov St 14 A</li>
        <li>Almaty, Almaty 050000</li>
        <li>Kazakhstan</li>
        <li>Phone: (+1)8 800 100 22 47</li>
        <li>Fax: (+1)8 800 100 22 47</li>
        <li>support@hexlet.io</li>
      </ul>
    </>
  );
}
