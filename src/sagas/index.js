import 'regenerator-runtime/runtime';
import { call, put, takeLatest, select } from 'redux-saga/effects';
import actions from 'actions';
import History from 'AppHistory';
import {
  ORDER_PARAM,
  ASCENDING_ORDER_PARAM,
  AGGREGAT_PARAM,
  AGGREGAT_SIZE_PARAM,
} from 'components/Funds/Constants';
import Config from 'services/Config';
import Funds from 'services/Funds';

/**
 * Saga of for retrieving config
 * @yield {Function} Saga effects to sequence flow of work
 */
export function* initSaga() {
  try {
    const config = yield call(Config.getConfig);
    yield put(actions.getConfigSucceeded(config));
  } catch (e) {
    yield put(actions.getConfigFailed(e));
  }
}

/**
 * Saga of getFunds action
 * @yield {Function} Saga effects to sequence flow of work
 */
export function* getFundsSaga() {
  try {
    const funds = yield call(Funds.getFunds);
    yield put(actions.getFundsSucceeded(funds));
  } catch (e) {
    yield put(actions.getFundsFailed(e));
  }
}

/**
 * Saga of updating url from filter/agregate
 * @yield {Function} Saga effects to sequence flow of work
 */
export function* updateUrlSaga() {
  const { filters, order, aggregat } = yield select(state => state.funds);

  const params = Object.entries(filters)
    .filter(([, value]) => Boolean(value))
    .map(([key, value]) => `${key}=${encodeURIComponent(value)}`);

  if (order.key) {
    params.push(`${ORDER_PARAM}=${encodeURIComponent(order.key)}`);

    if (!order.descending) {
      params.push(ASCENDING_ORDER_PARAM);
    }
  }

  if (aggregat.key) {
    params.push(`${AGGREGAT_PARAM}=${encodeURIComponent(aggregat.key)}`);
    params.push(`${AGGREGAT_SIZE_PARAM}=${encodeURIComponent(aggregat.size)}`);
  }

  let query = '';
  if (params.length) {
    query = `?${params.join('&')}`;
  }

  yield call(History.push, `/${query}`);
}

/**
 * Sagas of app.
 * @yield {Function} Sagas
 */
export default function* appSaga() {
  yield takeLatest(actions.INIT, initSaga);
  yield takeLatest(actions.GET_FUNDS_REQUEST, getFundsSaga);
  yield takeLatest([actions.SET_FILTER, actions.SET_ORDER, actions.SET_AGGREGAT], updateUrlSaga);
}
