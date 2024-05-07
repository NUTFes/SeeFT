import { Place } from '@type/common';

export const post = async (url: string, data: Place) => {
  const place = data.place;
  const remark = data.remark;
  const postUrl = url +
    '?place=' + place +
    '&remark=' + remark;
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

export const put = async (url: string, data: Place) => {
  const place = data.place;
  const remark = data.remark;
  const putUrl = url +
    '?place=' + place +
    '&remark=' + remark;
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

export const destroy = async (url: string, data: Place) => {
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