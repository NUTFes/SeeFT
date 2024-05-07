import React, { useState } from 'react';
import { useRouter } from 'next/router';

import { Place } from "@type/common";
import Input from '@components/common/Input';
import { post } from '@api/place';
import InformationPageLayout from '@components/layout/InformationPageLayout';

export default function Places() {
  const router = useRouter();

  const [formData, setFormData] = useState<Place>({
    id: 0,
    place: '',
    remark: '',
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
    <InformationPageLayout title='集合場所登録' submitText='登録' onClick={() => { addPlaceInformation(formData); }}>
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