import React from 'react';
import GraphiQL from 'graphiql';

import { connection, mutate } from 'thunder-react';

import '../node_modules/graphiql/graphiql.css';

function graphQLFetcher({query, variables}: any) {
  if (query.startsWith("mutation")) {
    return mutate({
      query,
      variables,
    });
  }
  return {
    subscribe(subscriber: any) {
      const next = subscriber.next || subscriber;

      const subscription = connection.subscribe({
        query: query,
        variables: {},
        observer: ({state, valid, error, value}: any) => {
          if (valid) {
            next({data: value});
          } else {
            next({state, error});
          }
        }
      });

      return {
        unsubscribe() {
          return subscription.close();
        }
      };
    }
  };
}

export function GraphiQLWithFetcher() {
  return <GraphiQL fetcher={graphQLFetcher} />;
}
