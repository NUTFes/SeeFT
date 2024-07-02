import { atom } from 'recoil';
import { recoilPersist } from 'recoil-persist';

import { STORE_KEYS } from '@/store/storeKey';
import { User } from '@/type/common';

const { persistAtom } = recoilPersist();

export const userAtom = atom<User>({
  key: STORE_KEYS.USER_STATE,
  default: {
    id: 0,
    name: '',
    mail: '',
    gradeID: 0,
    departmentID: 0,
    bureauID: 0,
    roleID: 0,
    studentNumber: 0,
    tel: '',
    password: '',
    createdAt: '',
    updatedAt: '',
  },
  effects_UNSTABLE: [persistAtom],
});
