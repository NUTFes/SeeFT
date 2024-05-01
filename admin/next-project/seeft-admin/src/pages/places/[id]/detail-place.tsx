import clsx from 'clsx';
import Head from 'next/head';

import { get } from '@api/api_methods';
import { Place } from "@type/common";
import MainLayout from '@components/layout/MainLayout';
import Button from '@components/common/Button';
import Input from '@components/common/Input';
import React, { useState } from 'react';

import { post } from '@api/place';
import { useRouter } from 'next/router';

interface Props {
  place: Place;
}

export const getServerSideProps = async (
  { params }: { params: { id: string } }) => {
  const placeID = params.id;
  const getPlaceURL = process.env.SSR_API_URI + '/places/' + placeID;
  const placeRes = await get(getPlaceURL);

  return {
    props: {
      place: placeRes,
    },
  };
};

export default function Places(props: Props) {
  const { place } = props;
  const router = useRouter();

  const [formData, setFormData] = useState<Place>({
    id: place.id,
    place: place.place,
    remark: place.remark,
  });

  const handler = (input: string) =>
    (e: React.ChangeEvent<HTMLSelectElement> | React.ChangeEvent<HTMLInputElement>) => {
      setFormData({ ...formData, [input]: e.target.value });
    }

  const addPlaceInformation = async (data: Place) => {
    const addPlaceInformationUrl = process.env.CSR_API_URI + '/places';
    await post(addPlaceInformationUrl, data);
    router.push('/places');
  };

  return (
    <MainLayout>
      <div className='mx-auto relative md:w-1/2 h-full bg-white-0 p-8'>
        <div className=''>
          <div className='mx-auto w-fit text-xl text-emphasis mb-8'>
            集合場所詳細
          </div>
          <div className='flex flex-col gap-3'>
            <div className='my-4 flex flex-col items-center justify-items-end gap-5 text-base text-emphasis'>
              <div className='flex w-full items-center'>
                <div className='flex w-1/4'>集合場所</div>
                <div className='col-span-4 w-full'>
                  <Input className='w-full' value={formData.place} onChange={handler('place')} />
                </div>
              </div>
              <div className='flex w-full items-center'>
                <div className='flex w-1/4'>備考</div>
                <div className='col-span-4 w-full'>
                  <Input className='w-full' value={formData.remark} onChange={handler('remark')} />
                </div>
              </div>
            </div>
          </div>
          <div className='mx-auto w-fit text-emphasis mb-8'>
            <Button className='bg-surface-2 hover:bg-surface-1'
              onClick={() => {
                addPlaceInformation(formData);
              }}>
              登録
            </Button>
          </div>
        </div>
      </div >
    </MainLayout >
  );
}