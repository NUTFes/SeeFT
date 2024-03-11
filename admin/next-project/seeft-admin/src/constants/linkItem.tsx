import { ReactNode } from 'react';
import { BsListTask } from 'react-icons/bs';
import { FaMapMarkedAlt } from 'react-icons/fa';
import { MdOutlineSchedule, MdManageAccounts } from "react-icons/md";

interface LinkItemProps {
  name: string;
  icon: ReactNode;
  href: string;
}

export const SeeftLinkItems: LinkItemProps[] = [
  {
    name: 'シフト',
    icon: <MdOutlineSchedule />,
    href: '/seefts'
  },
  {
    name: 'ユーザー',
    icon: <MdManageAccounts />,
    href: '/users'
  },
  {
    name: 'タスク',
    icon: <BsListTask />,
    href: '/tasks'
  },
  {
    name: '集合場所',
    icon: <FaMapMarkedAlt />,
    href: '/places'
  }
];
