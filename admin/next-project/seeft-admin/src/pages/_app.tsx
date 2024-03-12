import type { AppProps } from 'next/app';
import Head from 'next/head';
import { RecoilRoot } from 'recoil';

import 'tailwindcss/tailwind.css';

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <RecoilRoot>
      <Head>
        <link rel='icon' href='/icon.ico' />
      </Head>
      <div>
        <Component {...pageProps} />
      </div>
    </RecoilRoot>
  );
}

export default MyApp;