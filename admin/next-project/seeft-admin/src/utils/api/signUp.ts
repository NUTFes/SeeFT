import { User } from '@type/common';

export const signUp = async (url: string, data: User) => {
  const name = data.name;
  const mail = data.mail;
  const gradeID = data.gradeID;
  const bureauID = data.bureauID;
  const departmentID = data.departmentID;
  const roleID = data.roleID;
  const studentNumber = data.studentNumber;
  const tel = data.tel;
  const password = data.password
  const postUrl = url +
    '?name=' + name +
    '&mail=' + mail +
    '&student_number=' + studentNumber +
    '&grade_id=' + gradeID +
    '&department_id=' + departmentID +
    '&bureau_id=' + bureauID +
    '&role_id=' + roleID +
    '&tel=' + tel +
    '&password=' + password;
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