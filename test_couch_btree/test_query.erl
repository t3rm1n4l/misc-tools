-module(test_query).
-include("couch_db.hrl").
-export([start/0]).

calln(_Fn, 0) ->
    ok;
calln(Fn, N) ->
    Start = os:timestamp(),
    Fn(),
    io:format("time ~p~n", [timer:now_diff(os:timestamp(), Start)]),
    calln(Fn, N-1).


start() -> 
    Filename = "/Users/sarath/development/couchbase/ns_server/data/n_0/data/@indexes/default/main_da1eaf6fac28abafd16daa38c3bbbfd7.view.1",
    couch_file_write_guard:sup_start_link(),
    {ok, Fd}=couch_file:open(Filename),
    Root = {41098790,<<0,0,15,66,63,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255,255>>,18982105},    %{ok, Btree}=couch_btree:open(NR, Fd),
    NewBtree0 = #btree{root=Root, fd=Fd, binary_mode=true},
    NewBtree = couch_btree:set_options(NewBtree0, [{less, less_fun}]),

    WrapperFn = fun() ->
    {_,_,Val}  = couch_btree:foldl(NewBtree, fun(KV, Red, X) ->
                                                 case X of
                                                 10 ->
                                                     {stop, X};
                                                 _ ->
                                                     {ok, X+1}
                                                 end
                                             end, 0),
    Val
    end,
    Ops = 5500,
    Start = os:timestamp(),
    calln(WrapperFn, Ops),
    io:format("~p ops/secs~n", [Ops/(timer:now_diff(os:timestamp(), Start)/1000000)]).
