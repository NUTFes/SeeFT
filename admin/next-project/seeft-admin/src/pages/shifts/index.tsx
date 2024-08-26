import clsx from 'clsx';
import React, { use, useEffect, useMemo, useState } from 'react';

import { get } from '@api/api_methods';
import { User, Task, Bureau, Time, Shift, Weather, Date } from "@type/common";
import { post, put, destroy } from '@api/shift';
import { Select, ToggleButton } from '@components/common';
import ListPageLayout from '@components/layout/ListPageLayout';
import { TimeItems, TimeScaleItem } from '@constants/timeItem';
import { WeatherItem } from '@constants/weatherItem';
import { YearItem } from '@constants/yearItem';
import { DateItem } from '@constants/dateItem';
import { FaChevronRight } from "react-icons/fa";
import { FaChevronLeft } from "react-icons/fa";
import ReactSelect from 'react-select';

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
  const [timeScaleID, setTimeScaleID] = useState(3);
  const [hasDestoryMode, setHasDestoryMode] = useState(false);
  const [filteredBureau, setFilteredBureau] = useState<number>(0);
  const [selectedBureau, setSelectedBureau] = useState<number>(0);

  const [formData, setFormData] = useState<Shift>({
    id: 0,
    taskID: tasks[0].id,
    userID: 0,
    yearID: YearItem[YearItem.length - 1].id,
    dateID: 2,
    timeID: 0,
    weatherID: 1,
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
      setFormData({ ...formData, [input]: Number(e.target.value) });
    };

  const bureauHandler = () =>
    (e: React.ChangeEvent<HTMLSelectElement>) => {
      setSelectedBureau(Number(e.target.value));
    };

  const filterBureauHandler = () =>
    (e: React.ChangeEvent<HTMLSelectElement>) => {
      setFilteredBureau(Number(e.target.value));
    };

  const filteredTasks = useMemo(() => {
    return selectedBureau === 0 ? tasks
      : tasks.filter((task: Task) => (
        task.bureauID === selectedBureau
      ));
  }, [selectedBureau])

  const filteredUsers = useMemo(() => {
    return filteredBureau === 0 ? users.sort((a: User, b: User) => (a.bureauID - b.bureauID))
      : users.filter((user: User) => (
        user.bureauID === filteredBureau
      ))
  }, [filteredBureau]);

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
  };

  const timeList = useMemo(() => {
    return TimeItems.slice(((timeScaleID - 1) * 16), (timeScaleID * 16));
  }, [timeScaleID]);

  useEffect(() => {
    setFilteredShifts(shifts.filter((shift: Shift) => (
      shift.yearID === formData['yearID']
      && shift.dateID === formData['dateID']
      && shift.weatherID === formData['weatherID']
    )));
  }, [formData['yearID'], formData['dateID'], formData['weatherID'], shifts]);
  
  // タスクをマルチセレクタで扱えるように変換
  const formattedTasks = useMemo(() => { 
    return filteredTasks.map((task: Task) => ({
      value: task.id,
      label: task.task
    }))}, [filteredTasks, selectedBureau]);
  // セレクトされたタスクをstateで管理
  const [selectedTasks, setSelectedTasks] = useState(
    formattedTasks
  );
  useEffect(() => {
    setSelectedTasks(formattedTasks);
  }, [formattedTasks]);
  // タスクの選択状態を更新する関数
  const selectedTasksHandler = (selectedTasks: any) => {
    setSelectedTasks(selectedTasks.sort((a: any, b: any) => a.value - b.value));
  }

  return (
    <ListPageLayout title='シフト一覧'>
      <div className='my-3 border border-accent-1'>
        <div className='w-full flex justify-center items-center gap-6 p-1 '>
          <div className='w-1/6'>
            <Select className="w-full" value={filteredBureau} onChange={filterBureauHandler()}>
              <option key={0} value={0}>全局</option>
              {bureaus.map((data) => (
                <option key={data.id} value={data.id}>
                  {data.bureau}
                </option>
              ))}
            </Select>
          </div>
          <div className='w-1/6'>
            <Select className="w-full" value={formData.dateID} onChange={handler('dateID')}>
              {DateItem.map((data) => (
                <option key={data.id} value={data.id}>
                  {data.date}
                </option>
              ))}
            </Select>
          </div>
          <div className='w-1/6'>
            <Select className="w-full" value={formData.weatherID} onChange={handler('weatherID')}>
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
              {filteredUsers ? filteredUsers.map((user: User, index) => (
                <tr key={user.id}>
                  <td
                    className={clsx(
                      'px-1 py-1 bg-surface-2',
                      index === 0 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                      index === filteredUsers.length - 1 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{bureaus.find((bureau: Bureau) => (bureau.id === user.bureauID))?.bureau}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-1 bg-surface-2 border-accent-1 ',
                      index === 0 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                      index === filteredUsers.length - 1 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
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
      <div className='flex gap-4'>
        <div className='object-bottom w-1/3 border border-accent-1 my-2'>
          <div>
            シフトの編集方法
          </div>
          <div className='pl-4'>
            <div>
              <p>1. 担当局を選択</p>
              <div className='flex justify-center items-center gap-4 pl-4'>
                <div className='w-1/3'>担当局</div>
                <Select className='w-full' value={selectedBureau} onChange={bureauHandler()}>
                  <option key={0} value={0}>All</option>
                  {bureaus.map((data) => (
                    <option key={data.id} value={data.id}>
                      {data.bureau}
                    </option>
                  ))}
                </Select>
              </div>
            </div>
            <div>
              <p>2. 入力するシフトを選択</p>
              <div className='flex justify-center items-center gap-4 pl-4'>
                <div className='w-1/3'>シフト検索</div>
                <Select className='w-full' value={formData.taskID} onChange={handler('taskID')}>
                  {filteredTasks.length > 0 ? filteredTasks.map((data) => (
                    <option key={data.id} value={data.id}>
                      {data.task}
                    </option>
                  )) : <option key={1} value={1}>シフトが見つかりません</option>}
                </Select>
              </div>
            </div>
            <div>
              <p>3. 開始時刻のセルをクリック&ドラッグ</p>
            </div>
            <div className='pt-2'>
              <p>* [編集モード]/[削除モード]をクリックで切り替えられます</p>
            </div>
          </div>
        </div>
        <div className='object-bottom w-2/3 border border-accent-1 my-2'>
          <div>
            シフトの現在の人数
          </div>
          <div className='flex justify-center items-center gap-4 pl-4'>
            <div className='w-1/3'>表示するシフトを選択</div>
            <ReactSelect className='w-full z-20'
              closeMenuOnSelect={false}
              value={selectedTasks}
              defaultValue={formattedTasks}
              isMulti
              options={formattedTasks}
              isSearchable
              noOptionsMessage={() => 'シフトが見つかりません'}
              placeholder='シフトを選択'
              styles={{
                multiValueLabel: (provided) => ({
                    ...provided,
                    minWidth: '40px',  // 各選択されたoptionのラベルの最小幅を指定
                    maxWidth: '40px', // 各選択されたoptionのラベルの最大幅を指定
                    minHeight: '26.4px', // 各選択されたoptionのラベルの最小高さを指定
                    overflow: 'hidden', // 幅を超えた場合の処理
                    textOverflow: 'ellipsis',  // 溢れたテキストに省略記号を追加
                    whiteSpace: 'nowrap', // テキストを折り返さないように設定
                  }),
                menu: (provided) => ({
                  ...provided,
                  overflowY: 'auto',
                }),
                // ドロップダウンの高さを指定
                menuList: (provided) => ({
                  ...provided,
                  maxHeight: (selectedTasks.length) < 5 ? '150px' : '200px',
                }),
                option: (provided, state) => ({
                  ...provided,
                  minHeight: '40px',
                  color: state.isSelected ? 'white' : 'black',
                  backgroundColor: state.isSelected ? '#3182ce' : 'white',
                  '&:hover': {
                    backgroundColor: '#3182ce',
                    color: 'white',
                  },
                }),
              }}
              onChange={selectedTasksHandler}
            />
          </div>
          <div className='max-h-64 px-2 pb-2 overflow-y-auto select-none'>
            <table className='table-fixed mb-5 w-full border-collapse'>
              <thead className='sticky top-0 z-10'>
                <tr>
                  <th className='w-1/12 bg-surface-2 border border-accent-1 py-1'>
                    <p className='text-center text-sm text-emphasis'>タスク名</p>
                  </th>
                  <th className='w-1/12 bg-surface-2 border border-accent-1 py-1'>
                    <p className='text-center text-sm text-emphasis'>最大人数</p>
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
                {tasks ? tasks
                  .filter((task: Task) => (task.id === selectedTasks.find((option) => option.label === task.task)?.value))
                  .sort((a: Task, b: Task) => (a.id - b.id))
                  .map((task: Task, index) => (
                  <tr key={task.id}>
                    <td
                      className={clsx(
                        'px-1 py-1 bg-surface-2',
                        index === 0 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                        index === filteredUsers.length - 1 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                      )}
                    >
                      <p className='text-center text-sm text-emphasis truncate'>{selectedTasks[index].label}</p>
                    </td>
                    <td
                      className={clsx(
                        'px-1 py-1 bg-surface-2 border-accent-1 ',
                        index === 0 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                        index === filteredUsers.length - 1 ? 'pb-2 pt-1' : 'border border-accent-1 py-1',
                      )}
                    >
                      <p className='text-center text-sm text-emphasis'>{tasks[index].maxMember}</p>
                    </td>
                    {timeList.map((time: Time, i: number) => {
                      const currentMemberCount = filteredShifts
                            .filter((shift: Shift) => (shift.taskID === task.id))
                            .filter((shift: Shift) => (shift.dateID === formData.dateID))
                            .filter((shift: Shift) => (shift.timeID === time.id))
                            .length
                      const excessMemberColor = '#ffaaaa';
                      const shortageMemberColor = '#aaccff';
                      const backgroundColor = currentMemberCount < task.maxMember ? shortageMemberColor : currentMemberCount > task.maxMember ? excessMemberColor : '#ffffff';
                      return (
                        <td className='fixed-width w-3/64 bg-white-0 border border-accent-1 py-1 overflow-hidden text-ellipsis whitespace-nowrap text-center text-sm text-emphasis'
                          style={{background: (backgroundColor)}}
                        >
                          {currentMemberCount}
                        </td>
                      );
                    })}
                  </tr>
                )
                
                )
                 :
                  'ユーザーが登録されていません'
                }
              </tbody>
            </table>
          </div>
          <div className='flex justify-center items-center gap-10 my-4'>
            <div className='flex ml-auto '>
              <div  style={{width: 24, height: 24, background: '#aaccff'}}></div>
              <p>: 人数不足</p>
            </div>
            <div className='flex mr-6'>
              <div style={{width: 24, height: 24, background: '#ffaaaa'}}></div>
              <p>: 人数超過</p>
            </div>
          </div>
        </div>
      </div>
    </ListPageLayout >
  );
}