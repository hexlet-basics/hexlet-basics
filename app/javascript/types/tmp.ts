// remove when fixed https://github.com/skryukov/typelizer/issues/52
import type {
  LanguageCategoryCrud,
  LanguageCategoryItemCrud,
  LanguageCategoryQnaItemCrud,
  LanguageLandingPageCrud,
  LanguageLandingPageQnaItemCrud,
  SurveyCrud,
  SurveyItemCrud,
  SurveyScenarioCrud,
  SurveyScenarioItemCrud,
  SurveyScenarioTriggerCrud,
} from '@/types';

export type LanguageCategoryCrudDataWithAttrs = LanguageCategoryCrud['data'] & {
  items_attributes: Array<LanguageCategoryItemCrud>;
  qna_items_attributes: Array<LanguageCategoryQnaItemCrud>;
};

export type LanguageCategoryCrudWithAttrs = Omit<
  LanguageCategoryCrud,
  'data'
> & {
  data: LanguageCategoryCrudDataWithAttrs;
};

// Расширяем SurveyCrudData, добавляя items_attributes
export type SurveyCrudDataWithAttrs = SurveyCrud['data'] & {
  items_attributes: Array<SurveyItemCrud>;
};

// Создаём полный тип SurveyCrud с атрибутами
export type SurveyCrudWithAttrs = Omit<SurveyCrud, 'data'> & {
  data: SurveyCrudDataWithAttrs;
};

// Расширяем SurveyScenarioCrudData, добавляя *_attributes
export type SurveyScenarioCrudDataWithAttrs = SurveyScenarioCrud['data'] & {
  items_attributes: Array<SurveyScenarioItemCrud>;
  triggers_attributes: Array<SurveyScenarioTriggerCrud>;
};

// Создаём полный тип SurveyScenarioCrud с атрибутами
export type SurveyScenarioCrudWithAttrs = Omit<SurveyScenarioCrud, 'data'> & {
  data: SurveyScenarioCrudDataWithAttrs;
};

// Расширяем LanguageLandingPageCrudData, добавляя qna_items_attributes
export type LanguageLandingPageCrudDataWithAttrs =
  LanguageLandingPageCrud['data'] & {
    qna_items_attributes: Array<LanguageLandingPageQnaItemCrud>;
  };

// Создаём полный тип LanguageLandingPageCrud с атрибутами
export type LanguageLandingPageCrudWithAttrs = Omit<
  LanguageLandingPageCrud,
  'data'
> & {
  data: LanguageLandingPageCrudDataWithAttrs;
};
