import clsx from 'clsx';
import Head from 'next/head';

import { get } from '@api/api_methods';
import { User } from "@type/common";
import MainLayout from '@components/layout/MainLayout';

interface Props {
  users: User[];
}

export const getServerSideProps = async () => {
  const getUserURL = process.env.SSR_API_URI + '/users';
  const userRes = await get(getUserURL);

  return {
    props: {
      users: userRes,
    },
  };
};

export default function Uesrs(props: Props) {
  const { users } = props;

  return (
    <MainLayout>
      <div className='w-full h-full bg-surface-1/[0.5] flex-col p-8'>
        <div className='items-center'>
          test
        </div>
        <div>
          table
        </div>
      </div>
    </MainLayout >
  );
}