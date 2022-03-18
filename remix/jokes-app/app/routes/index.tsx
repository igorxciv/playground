import { LinksFunction, MetaFunction } from "remix";

import stylesUrl from '~/styles/index.css';

export default function Index() {
  return (
    <div>
      Hello World route!
    </div>
  );
}

export const links: LinksFunction = () => [
  {
    rel: 'stylesheet',
    href: stylesUrl,
  }
]

export const meta: MetaFunction = () => ({
  title: 'Remix: So great, its funny!',
  description: 'Remix jokes app. Learn Remix and laugh at the same time!',
})