import clsx from 'clsx';
import Head from 'next/head';
import React from 'react';

import Header from '@components/common/Header/Header';


interface LayoutProps {
  children?: React.ReactNode;
}

export default function MainLayout(props: LayoutProps) {
  return (
    <>
      <Head>
        <title className='text-emphasis'>SeeFT</title>
      </Head>
      <div className={clsx('h-screen w-full')}>
        <div className={clsx('h-10 w-full')}>
          <Header />
        </div>
        <div>
          {props.children}
        </div>
      </div>
    </>
  );
}