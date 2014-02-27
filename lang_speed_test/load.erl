-module(load).
-export([main/1]).

gen_items(N) ->
    gen_items(N, []).

gen_items(0, A) ->
    A;
gen_items(N, A) ->
    gen_items(N-1, [io_lib:format("doc_~p", [N])|A]).
    %gen_items(N-1, [binary_to_list(iolist_to_binary(io_lib:format("doc_~7..0B", [N])))|A]).

main(_) ->
    N = 500000000,
    Start=os:timestamp(),
    gen_items(N),
    Diff=timer:now_diff(os:timestamp(), Start)/1000000,
    io:format("took ~.3fs\n", [Diff]).

