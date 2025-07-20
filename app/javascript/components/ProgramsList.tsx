import { Carousel } from '@mantine/carousel';
import { useMantineTheme } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import type { LanguageLandingPageForLists } from '@/types';
import CourseBlock from './CourseBlock';

type Props = {
  landingPages: LanguageLandingPageForLists[];
  attributes?: Record<string, string>;
};

export default function CoursesList({ landingPages, attributes }: Props) {
  if (landingPages.length === 0) {
    return '<!-- CoursesList is empty -->';
  }
  return (
    <Carousel
      slideSize={{ base: '50%', xs: '50%' }}
      slideGap={{ base: 'xs', xs: 'xl' }}
      emblaOptions={{
        loop: true,
        dragFree: false,
        align: 'center',
        // slidesToScroll: mobile ? 1 : 2,
      }}
    >
      {landingPages.map((lp) => (
        <Carousel.Slide key={lp.id}>
          <CourseBlock landingPage={lp} />
        </Carousel.Slide>
      ))}
    </Carousel>
  );
}
