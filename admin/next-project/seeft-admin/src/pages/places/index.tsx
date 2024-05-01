import clsx from 'clsx';
import { useRouter } from 'next/router';

import { get } from '@api/api_methods';
import { Place } from "@type/common";
import MainLayout from '@components/layout/MainLayout';
import { MdEdit, MdDeleteForever } from "react-icons/md";
import Button from '@components/common/Button';
import { destroy } from '@api/place';

interface Props {
  places: Place[];
}

export const getServerSideProps = async () => {
  const getPlaceURL = process.env.SSR_API_URI + '/places';
  const placeRes = await get(getPlaceURL);

  return {
    props: {
      places: placeRes,
    },
  };
};

export default function Uesrs(props: Props) {
  const { places } = props;
  const router = useRouter();

  const addPlacePageRouter = () => {
    router.push('places/add-place');
  }

  const PlaceDetailPageRouter = (place: Place) => {
    router.push('places/' + place.id + '/detail-place');
  }

  const destroyPlaceInformation = async (data: Place) => {
    const destroyPlaceInformationUrl = process.env.CSR_API_URI + '/places';
    await destroy(destroyPlaceInformationUrl, data);
    router.reload();
  };

  return (
    <MainLayout>
      <div className='w-full h-full bg-white-0 flex-col p-8'>
        <div className='items-center text-xl text-emphasis'>
          集合場所一覧
        </div>
        <div className='items-center'>
          <div className='text-right pr-4'>
            <Button className='bg-surface-2 border-accent-2 text-right text-emphasis pr-4 hover:bg-surface-1' onClick={addPlacePageRouter}>
              集合場所追加
            </Button>
          </div>
        </div>
        <div className='p-5'>
          <table className='mb-5 w-full table-auto border-collapse'>
            <thead>
              <tr>
                <th className='w-1/4 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                  <p className='text-center text-sm text-emphasis'>場所</p>
                </th>
                <th className='w-1/4 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3'>
                  <p className='text-center text-sm text-emphasis'>備考</p>
                </th>
                <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3' />
                <th className='w-1/12 border border-x-white-0 border-b-accent-1 border-t-white-0 py-3' />
              </tr>
            </thead>
            <tbody className='border border-x-white-0 border-b-accent-1 border-t-white-0'>
              {places ? places.map((place: Place, index) => (
                <tr key={place.id}>
                  <td
                    className={clsx(
                      'px-1 py-2',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === places.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{place.place}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-2',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === places.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                    )}
                  >
                    <p className='text-center text-sm text-emphasis'>{place.remark}</p>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-2 rounded-full hover:bg-accent-1',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === places.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                    )}
                  >
                    <div className='flex justify-items-center gap-4'
                      onClick={() => { PlaceDetailPageRouter(place) }}>
                      <MdEdit />
                      <p className='text-center text-sm text-emphasis'>
                        編集
                      </p>
                    </div>
                  </td>
                  <td
                    className={clsx(
                      'px-1 py-2 rounded-full hover:bg-accent-1',
                      index === 0 ? 'pb-3 pt-4' : 'py-3',
                      index === places.length - 1 ? 'pb-4 pt-3' : 'border-b-accent-1 py-3',
                    )}
                  >
                    <div className='flex justify-items-center gap-4'
                      onClick={() => { destroyPlaceInformation(place); }}>
                      <MdDeleteForever />
                      <p className='text-center text-sm text-emphasis'>
                        削除
                      </p>
                    </div>
                  </td>
                </tr>
              )) : null}
            </tbody>
          </table>
        </div>
      </div>
    </MainLayout >
  );
}