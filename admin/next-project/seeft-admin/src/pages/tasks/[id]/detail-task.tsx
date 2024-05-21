import React, { useState } from 'react';
import { useRouter } from 'next/router';

import { get } from '@api/api_methods';
import { Task, User, Place } from "@type/common";
import { Input, Select } from '@components/common';
import { put } from '@api/task';
import InformationPageLayout from '@components/layout/InformationPageLayout';

interface Props {
  task: Task;
  users: User[];
  places: Place[];
}

export const getServerSideProps = async (
  { params }: { params: { id: string } }
) => {
  const taskID = params.id
  const getTaskURL = process.env.SSR_API_URI + '/tasks/' + taskID;
  const getUserURL = process.env.SSR_API_URI + '/users';
  const getPlaceURL = process.env.SSR_API_URI + '/places';
  const taskRes = await get(getTaskURL);
  const userRes = await get(getUserURL);
  const placeRes = await get(getPlaceURL);

  return {
    props: {
      task: taskRes,
      users: userRes,
      places: placeRes,
    },
  };
};

export default function Users(props: Props) {
  const { task, users, places } = props;
  const router = useRouter();

  const [formData, setFormData] = useState<Task>({
    id: task.id,
    task: task.task,
    placeID: task.placeID,
    url: task.url,
    superviserID: task.superviserID,
    color: task.color,
    remark: task.remark,
    yearID: task.yearID,
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

  const updateTaskInformation = async (data: Task) => {
    const putTaskInformationUrl = process.env.CSR_API_URI + '/tasks/' + data.id;
    await put(putTaskInformationUrl, data);
    router.push('/tasks');
  };

  return (
    <InformationPageLayout title='タスク登録' submitText='登録' onClick={() => { updateTaskInformation(formData); }}>
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