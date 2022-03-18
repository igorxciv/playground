import { Links, LinksFunction, LiveReload, Meta, MetaFunction, Outlet, Scripts } from "remix";

import globalStylesUrl from './styles/global.css';
import globalMediumStylesUrl from './styles/global-medium.css';
import globalLargeStylesUrl from './styles/global-large.css';

function Document({ children, title = 'Remix: So great, its funny!' }: { children?: React.ReactNode, title?: string }) {
  return (
    <html lang="en">
      <head>
        <meta charSet="utf-8" />
        <Meta />
        <title>{title}</title>
        <Links />
      </head>
      <body>
        {children}
        <Scripts />
        <LiveReload />
      </body>
    </html>
  )
}

export default function App() {
  return (
    <Document>
      <Outlet />
    </Document>
  );
}

export const meta: MetaFunction = () => {
  const description = 'Learn Remix and laugh at the same time!';
  return {
    description,
    keywords: 'Remix, jokes',
    'twitter:image': 'https://remix-jokes.lol/social.png',
    "twitter:card": "summary_large_image",
    "twitter:creator": "@remix_run",
    "twitter:site": "@remix_run",
    "twitter:title": "Remix Jokes",
    "twitter:description": description,
  }
}

export const links: LinksFunction = () => [
  {
    rel: "stylesheet",
    href: globalStylesUrl,
  },
  {
    rel: 'stylesheet',
    href: globalMediumStylesUrl,
    media: 'print, (min-width: 640px)',
  },
  {
    rel: 'stylesheet',
    href: globalLargeStylesUrl,
    media: 'screen and (min-width: 1024px)',
  }
]

export function ErrorBoundary({ error }: { error: Error }) {
  return (
    <Document title="Uh-oh!">
      <div className="error-container">
        <h1>App error!</h1>
        <pre>{error.message}</pre>
      </div>
    </Document>
  )
}
