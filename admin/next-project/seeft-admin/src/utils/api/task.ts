import { Task } from '@type/common';

export const post = async (url: string, data: Task) => {
  const task = data.task;
  const placeID = data.placeID;
  const manual = data.url;
  const bureauID = data.bureauID;
  const maxMember = data.maxMember;
  const color = data.color;
  const remark = data.remark;
  const yearID = data.yearID;
  const postUrl = url +
    '?task=' + task +
    '&place_id=' + placeID +
    '&url=' + manual +
    '&bureau_id=' + bureauID +
    '&max_member=' + maxMember +
    '&color=' + color +
    '&remark=' + remark +
    '&year_id=' + yearID;
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

export const put = async (url: string, data: Task) => {
  const task = data.task;
  const placeID = data.placeID;
  const manual = data.url;
  const bureauID = data.bureauID;
  const maxMember = data.maxMember;
  const color = data.color;
  const remark = data.remark;
  const yearID = data.yearID;
  const putUrl = url +
    '?task=' + task +
    '&place_id=' + placeID +
    '&url=' + manual +
    '&bureau_id=' + bureauID +
    '&max_member=' + maxMember +
    '&color=' + color +
    '&remark=' + remark +
    '&year_id=' + yearID;
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

export const destroy = async (url: string, data: Task) => {
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