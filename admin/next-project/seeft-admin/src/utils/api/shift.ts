import { Shift } from '@type/common';

export const post = async (url: string, data: Shift) => {
  const taskID = data.taskID;
  const userID = data.userID;
  const yearID = data.yearID;
  const dateID = data.dateID;
  const timeID = data.timeID;
  const weatherID = data.weatherID;
  const isAttendance = data.isAttendance;
  const postUrl = url +
    '?task_id=' + taskID +
    '&user_id=' + userID +
    '&year_id=' + yearID +
    '&date_id=' + dateID +
    '&time_id=' + timeID +
    '&weather_id=' + weatherID +
    '&is_attendance=' + isAttendance;
  const res = await fetch(postUrl, {
    method: 'POST',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  return await res;
  // return await res.json();
};

export const put = async (url: string, data: Shift) => {
  const taskID = data.taskID;
  const userID = data.userID;
  const yearID = data.yearID;
  const dateID = data.dateID;
  const timeID = data.timeID;
  const weatherID = data.weatherID;
  const isAttendance = data.isAttendance;
  const putUrl = url +
    '?task_id=' + taskID +
    '&user_id=' + userID +
    '&year_id=' + yearID +
    '&date_id=' + dateID +
    '&time_id=' + timeID +
    '&weather_id=' + weatherID +
    '&is_attendance=' + isAttendance;
  const res = await fetch(putUrl, {
    method: 'PUT',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  return await res.json();
};

export const destroy = async (url: string, data: Shift) => {
  const id = data.id
  const destroyUrl = url + '?id=' + id;
  const res = await fetch(destroyUrl, {
    method: 'DELETE',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  return await res;
};