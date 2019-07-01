import actions from 'actions';
import reducer, { initialState } from './funds';

it('should return initial state', () => {
  expect(reducer(undefined, { type: '' })).toEqual(initialState);
});

it('should update with given funds on fetch succeed', () => {
  const funds = [{ id: 8000 }];

  expect(reducer(initialState, actions.getFundsSucceeded(funds))).toEqual({
    ...initialState,
    funds,
  });
});

it('should remove funds without id', () => {
  const funds = [{ id: 8000 }, { name: 'test' }];

  expect(reducer(initialState, actions.getFundsSucceeded(funds))).toEqual({
    ...initialState,
    funds: [{ id: 8000 }],
  });
});
