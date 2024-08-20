import { SignIn } from '@type/common';

export const signIn = async (url: string, data: SignIn) => {
  const studentNumber = data.studentNumber;
  const password = data.password;
  const postUrl = url + '?student_number=' + studentNumber + '&password=' + password;
  const res = await fetch(postUrl, {
    method: 'POST',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
  });
  return await res;
};