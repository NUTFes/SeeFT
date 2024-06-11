import clsx from 'clsx';
import React, { useEffect, useMemo, useState } from 'react';
import { useRouter } from 'next/navigation';

import { get } from '@api/api_methods';
import { User, Task, Bureau, Time, Shift, Weather, Date } from "@type/common";
import MainLayout from '@components/layout/MainLayout';
import Button from '@components/common/Button';
import { post, put, destroy } from '@api/shift';
import { Select, ToggleButton } from '@components/common';
import ListPageLayout from '@components/layout/ListPageLayout';
import { TimeItems, TimeScaleItem } from '@constants/timeItem';
import { WeatherItem } from '@constants/weatherItem';
import { YearItem } from '@constants/yearItem';
import { DateItem } from '@constants/dateItem';
import { FaChevronRight } from "react-icons/fa";
import { FaChevronLeft } from "react-icons/fa";

interface Props {
  shifts: Shift[];
  users: User[];
  tasks: Task[];
  bureaus: Bureau[];
}

export const getServerSideProps = async () => {
  const getShiftURL = process.env.SSR_API_URI + '/shifts-admin';
  const getUserURL = process.env.SSR_API_URI + '/users';
  const getTaskURL = process.env.SSR_API_URI + '/tasks';
  const getBureauURL = process.env.SSR_API_URI + '/bureaus';
  const shiftRes = await get(getShiftURL);
  const userRes = await get(getUserURL);
  const taskRes = await get(getTaskURL);
  const bureauRes = await get(getBureauURL);

  return {
    props: {
      shifts: shiftRes,
      users: userRes,
      tasks: taskRes,
      bureaus: bureauRes,
    },
  };
};

export default function Users(props: Props) {
  const { users, tasks, bureaus } = props;
  const [shifts, setShifts] = useState(props.shifts);
  const [filteredShifts, setFilteredShifts] = useState<Shift[]>([]);
  const [isMouseDown, setIsMouseDown] = useState<boolean>(false);
  const [yearID, setYearID] = useState(YearItem[YearItem.length - 1].id);
  const [dateID, setDateID] = useState(2);
  const [timeScaleID, setTimeScaleID] = useState(3);
  const [weatherID, setWeatherID] = useState(1);
  const [hasDestoryMode, setHasDestoryMode] = useState(false);
  const router = useRouter();

  const [formData, setFormData] = useState<Shift>({
    id: 0,
    taskID: tasks[0].id,
    userID: 0,
    yearID: yearID,
    dateID: dateID,
    timeID: 0,
    weatherID: weatherID,
    isAttendance: false
  });

  // シフトの追加API(DBの仕様上使用していない)
  const addShiftInformation = async (data: Shift, user: User, time: Time) => {
    const addData = { ...data, id: shifts[shifts.length - 1].id + 1, userID: user.id, timeID: time.id }
    const addUserInformationUrl = process.env.CSR_API_URI + '/shifts-admin';
    await post(addUserInformationUrl, addData);
    const updatedShifts = [...shifts, addData];
    setShifts(updatedShifts);
  };

  // シフトの編集API
  const updateShiftInformation = async (data: Shift, user: User, time: Time, id: number) => {
    const updateData = { ...data, id: id, taskID: Number(data.taskID), userID: user.id, timeID: time.id };
    const putShiftInformationUrl = process.env.CSR_API_URI + '/shifts-admin/' + id;
    await put(putShiftInformationUrl, updateData);
    const updatedShifts = shifts.map((shift: Shift) => (shift.id === updateData.id ? updateData : shift));
    setShifts(updatedShifts);
  };

  // シフトの削除API
  const destroyShiftInformation = async (data: Shift, user: User, time: Time, id: number) => {
    const updateData = { ...data, id: id, taskID: 1, userID: user.id, timeID: time.id }
    const putShiftInformationUrl = process.env.CSR_API_URI + '/shifts-admin/' + id;
    await put(putShiftInformationUrl, updateData);
    const updatedShifts = shifts.map((shift: Shift) => (shift.id === updateData.id ? updateData : shift));
    setShifts(updatedShifts);
    // DBの仕様上destoryを実行しない方がよい(今後修正)
    // const destroyShiftInformationUrl = process.env.CSR_API_URI + '/shifts-admin';
    // await destroy(destroyShiftInformationUrl, data);
    // const updatedShifts = shifts.filter((shift: Shift) => (shift.id !== id));
  };

  const handler = (input: string) =>
    (e: React.ChangeEvent<HTMLSelectElement> | React.ChangeEvent<HTMLInputElement>) => {
      setFormData({ ...formData, [input]: e.target.value });
    };

  const handleMouseDown = (user: User, time: Time, id: number) => {
    setIsMouseDown(true);

    const formDataCopy = { ...formData }; // formDataのコピーを作成
    if (id) {
      if (hasDestoryMode) {
        destroyShiftInformation(formDataCopy, user, time, id);
      } else {
        updateShiftInformation(formDataCopy, user, time, id);
      }
    } else {
      addShiftInformation(formDataCopy, user, time);
    }
  };

  const handleMouseUp = () => {
    setIsMouseDown(false);
  };

  const handleMouseEnter = (user: User, time: Time, id: number) => {
    const formDataCopy = { ...formData };
    if (isMouseDown) {
      if (hasDestoryMode) {
        destroyShiftInformation(formDataCopy, user, time, id);
      } else {
        const existingShift = filteredShifts.find(shift => shift.userID === user.id && shift.timeID === time.id);
        if (existingShift) {
          updateShiftInformation(formDataCopy, user, time, existingShift.id);
        } else {
          addShiftInformation(formDataCopy, user, time);
        }
      }
    }
  };

  const dateHandler = () =>
    (e: React.ChangeEvent<HTMLSelectElement>) => {
      setDateID(Number(e.target.value));
    };

  const weatherHandler = () =>
    (e: React.ChangeEvent<HTMLSelectElement>) => {
      setWeatherID(Number(e.target.value));
    };

  const changeTimeScale = (value: number) => {
    if (value === 1) {
      if (timeScaleID === 6) {
        setTimeScaleID(1);
      }
      else {
        setTimeScaleID(timeScaleID + 1);
      }
    }
    else {
      if (timeScaleID === 1) {
        setTimeScaleID(6);
      }
      else {
        setTimeScaleID(timeScaleID - 1);
      }
    }
  }


  const timeList = useMemo(() => {
    return TimeItems.slice(((timeScaleID - 1) * 16), (timeScaleID * 16));
  }, [timeScaleID]);

  useEffect(() => {
    setFilteredShifts(shifts.filter((shift: Shift) => (
      shift.yearID === yearID
      && shift.dateID === dateID
      && shift.weatherID === weatherID
    )));
  }, [yearID, dateID, weatherID, shifts]);

  return (
    <ListPageLayout title='シフト一覧'>
      <div className='my-3 border border-accent-1'>
        <div className='w-full flex justify-center items-center gap-6 p-1 '>
          <div className='w-1/6'>
            <Select className="w-full" value={DateItem.find((date: Date) => (date.id === dateID))?.id} onChange={(dateHandler())}>
              {DateItem.map((data) => (
                <option key={data.id} value={data.id}>
                  {data.date}
                </option>
              ))}
            </Select>
          </div>
          <div className='w-1/6'>
            <Select className="w-full" value={WeatherItem.find((weather: Weather) => (weather.id === weatherID))?.id} onChange={(weatherHandler())}>
              {WeatherItem.map((data) => (
                <option key={data.id} value={data.id}>
                  {data.weather}
                </option>
              ))}
            </Select>
          </div>
          <div className='w-1/6'>
            <ToggleButton isToggleState={hasDestoryMode} initialLabel='編集モード' toggledLabel='削除モード' onClick={() => (setHasDestoryMode(!hasDestoryMode))} />
          </div>
        </div>
        <div className='flex justify-between items-center p-2'>
          <div className='flex justify-center items-center gap-4 px-2 rounded-lg border border-accent-1 hover:bg-accent-1' onClick={() => { changeTimeScale(-1) }}>
            <FaChevronLeft />
            {timeScaleID === 1 ? TimeScaleItem[5].time : TimeScaleItem[timeScaleID - 2].time}
          </div>
          <div className='text-xl'>
            {TimeScaleItem[timeScaleID - 1].time}
          </div>
          <div className='flex justify-center items-center gap-4 px-2 rounded-lg border border-accent-1 hover:bg-accent-1' onClick={() => { changeTimeScale(1) }}>
            {timeScaleID === 6 ? TimeScaleItem[0].time : TimeScaleItem[timeScaleID].time}
            <FaChevronRight />
          </div>
        </div>
        <div className='max-h-64 px-2 pb-2 overflow-y-auto select-none'>
          <table className='table-fixed mb-5 w-full border-collapse'>
            <thead className='sticky top-0 z-10'>
              <tr>
                <th className='w-1/12 bg-surface-2 border border-accent-1 py-1'>
                  <p className='text-center text-sm text-emphasis'>所属局</p>
                </th>
                <th className='w-2/12 bg-surface-2 border border-accent-1 py-1'>
                  <p className='text-center text-sm text-emphasis'>名前</p>
                </th>
                {timeList.map((time: Time, i: number) => (
                  i % 2 === 0 ?
                    <th className='w-3/64 bg-surface-2 border border-accent-1 py-1'>
                      <p className='text-center text-sm text-emphasis'>{time.time + '-'}</p>
                    </th>
                    :
                    <th className='w-3/64 bg-surface-2 border border-accent-1 py-1' />
                ))}
              </tr>
            </thead>
            <tbody className='border border-x-white-0 border-b-accent-1 border-t-white-0'>
              {users ? users.map((user: User, index) => (
                <tr key={user.id}>
                  <td
                    className={clsx(
                      'px-1 py-1 bg-surface-2',
                      index === 0 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                      index === users.length - 1 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{bureaus.find((bureau: Bureau) => (bureau.id === user.bureauID))?.bureau}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-1 bg-surface-2 border-accent-1 ',
                      index === 0 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                      index === users.length - 1 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{user.name}</p>
                  </td>
                  {timeList.map((time: Time, i: number) => {
                    const shift = filteredShifts.find((shift: Shift) => (shift.userID === user.id && shift.timeID === time.id));
                    const task = shift ? tasks.find(task => task.id === shift.taskID) : null;
                    const backgroundColor = task ? `#${task.color}` : '#ffffff';

                    return (
                      <td className='fixed-width w-3/64 bg-white-0 border border-accent-1 py-1 overflow-hidden text-ellipsis whitespace-nowrap text-center text-sm text-emphasis'
                        style={{
                          background: (backgroundColor)
                        }}
                        onMouseDown={() => handleMouseDown(user, time, shift?.id || 0)}
                        onMouseUp={handleMouseUp}
                        onMouseEnter={() => handleMouseEnter(user, time, shift?.id || 0)}
                      >
                        {tasks ? tasks.find((task: Task) => (task.id === (filteredShifts ? filteredShifts.find((shift: Shift) => (shift.userID === user.id && shift.timeID === time.id))?.taskID : null)))?.task
                          : null}
                      </td>
                    );
                  })}
                </tr>
              )) :
                'ユーザーが登録されていません'
              }
            </tbody>
          </table>
        </div>
      </div>
      <div className='object-bottom w-1/3 border border-accent-1 my-2'>
        <div>
          シフトの編集方法
        </div>
        <div className='pl-4'>
          <div>
            <p>1. 入力するシフトを選択</p>
            <div className='flex justify-center items-center gap-4 pl-4'>
              <div className='w-1/3'>シフト検索</div>
              <Select className='w-full' value={formData.taskID} onChange={handler('taskID')}>
                {tasks.map((data) => (
                  <option key={data.id} value={data.id}>
                    {data.task}
                  </option>
                ))}
              </Select>
            </div>
          </div>
          <div>
            <p>2. 開始時刻のセルをクリック&ドラッグ</p>
          </div>
          <div className='pt-2'>
            <p>* [編集モード]/[削除モード]をクリックで切り替えられます</p>
          </div>
        </div>
      </div>
    </ListPageLayout >
  );
}