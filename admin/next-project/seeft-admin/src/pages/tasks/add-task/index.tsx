import React, { useState } from 'react';
import { useRouter } from 'next/router';

import { get } from '@api/api_methods';
import { Task, User, Place } from "@type/common";
import { Input, Select } from '@components/common';
import { post } from '@api/task';
import InformationPageLayout from '@components/layout/InformationPageLayout';
import { url } from 'inspector';

interface Props {
  users: User[];
  places: Place[];
}

export const getServerSideProps = async () => {
  const getUserURL = process.env.SSR_API_URI + '/users';
  const getPlaceURL = process.env.SSR_API_URI + '/places';
  const userRes = await get(getUserURL);
  const placeRes = await get(getPlaceURL);

  return {
    props: {
      users: userRes,
      places: placeRes,
    },
  };
};

export default function Users(props: Props) {
  const { users, places } = props;
  const router = useRouter();

  const [formData, setFormData] = useState<Task>({
    id: 0,
    task: '',
    placeID: 1,
    url: '',
    superviserID: 1,
    color: 'fffafa',
    remark: '',
    yearID: 43,
  });

  const handler = (input: string) =>
    (e: React.ChangeEvent<HTMLSelectElement> | React.ChangeEvent<HTMLInputElement>) => {
      if (input === 'color') {
        setFormData({ ...formData, [input]: e.target.value.replace('#', '') });
      }
      else {
        setFormData({ ...formData, [input]: e.target.value });
      }
    }

  const addTaskInformation = async (data: Task) => {
    const addTaskInformationUrl = process.env.CSR_API_URI + '/tasks';
    await post(addTaskInformationUrl, data);
    router.push('/tasks');
  };

  return (
    <InformationPageLayout title='タスク登録' submitText='登録' onClick={() => { addTaskInformation(formData); }}>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>タスク名</div>
        <div className='col-span-4 w-full'>
          <Input className='w-full' value={formData.task} onChange={handler('task')} />
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4 whitespace-nowrap'>集合場所</div>
        <div className='col-span-4 w-full'>
          <Select className='w-full' value={formData.placeID} onChange={handler('placeID')}>
            {places.map((data) => (
              <option key={data.id} value={data.id}>
                {data.place}
              </option>
            ))}
          </Select>
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>マニュアルURL</div>
        <div className='col-span-4 w-full'>
          <Input className='w-full' value={formData.url} onChange={handler('url')} />
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>責任者</div>
        <div className='col-span-4 w-full'>
          <Select className='w-full' value={formData.superviserID} onChange={handler('superviserID')}>
            {users.map((data) => (
              <option key={data.id} value={data.id}>
                {data.name}
              </option>
            ))}
          </Select>
        </div>
      </div>
      <div className='flex w-full items-center'>
        <div className='flex w-1/4'>カラー</div>
        <div className='flex w-full items-center gap-4'>
          <div
            className="flex w-8 h-8 rounded-full border border-gray-400 "
            style={{ background: ('#' + formData.color) }}
          >
            <input
              type="color"
              value={'#' + formData.color}
              onChange={handler('color')}
              className="w-full h-full opacity-0"
            />
          </div>
          <Input className='flex w-full' value={formData.color} onChange={handler('color')} />
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