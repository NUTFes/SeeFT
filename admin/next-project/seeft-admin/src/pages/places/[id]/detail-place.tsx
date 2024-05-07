import React, { useState } from 'react';
import { useRouter } from 'next/router';

import { Place } from "@type/common";
import Input from '@components/common/Input';
import { get } from '@api/api_methods';
import { post } from '@api/place';
import InformationPageLayout from '@components/layout/InformationPageLayout';

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

  const updatePlaceInformation = async (data: Place) => {
    const putPlaceInformationUrl = process.env.CSR_API_URI + '/places';
    await post(putPlaceInformationUrl, data);
    router.push('/places');
  };

  return (
    <InformationPageLayout title='集合場所詳細' submitText='編集' onClick={() => { updatePlaceInformation(formData); }}>
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
    </InformationPageLayout>

  );
}